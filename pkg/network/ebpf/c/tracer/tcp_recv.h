#ifndef __TCP_RECV_H
#define __TCP_RECV_H

#include "bpf_helpers.h"
#include "bpf_telemetry.h"
#include "bpf_bypass.h"

#include "tracer/stats.h"
#include "tracer/events.h"
#include "tracer/maps.h"
#include "sock.h"

SEC("kprobe/tcp_recvmsg")
int BPF_KPROBE_INSTR(WITH(BYPASS), kprobe__tcp_recvmsg) {
    u64 pid_tgid = bpf_get_current_pid_tgid();
#if defined(COMPILE_RUNTIME) && LINUX_VERSION_CODE < KERNEL_VERSION(4, 1, 0)
    struct sock *skp = (struct sock*)PT_REGS_PARM2(ctx);
    int flags = (int)PT_REGS_PARM6(ctx);
#elif defined(COMPILE_RUNTIME) && LINUX_VERSION_CODE < KERNEL_VERSION(5, 19, 0)
    struct sock *skp = (struct sock*)PT_REGS_PARM1(ctx);
    int flags = (int)PT_REGS_PARM5(ctx);
#else
    struct sock *skp = (struct sock*)PT_REGS_PARM1(ctx);
    int flags = (int)PT_REGS_PARM4(ctx);
#endif
    if (flags & MSG_PEEK) {
        return 0;
    }

    bpf_map_update_with_telemetry(tcp_recvmsg_args, &pid_tgid, &skp, BPF_ANY);
    return 0;
}

#if defined(COMPILE_CORE) || defined(COMPILE_PREBUILT)

SEC("kprobe/tcp_recvmsg")
int BPF_KPROBE_INSTR(WITH(BYPASS), kprobe__tcp_recvmsg__pre_5_19_0) {
    u64 pid_tgid = bpf_get_current_pid_tgid();
    int flags = (int)PT_REGS_PARM5(ctx);
    if (flags & MSG_PEEK) {
        return 0;
    }
    struct sock *skp = (struct sock *)PT_REGS_PARM1(ctx);
    bpf_map_update_with_telemetry(tcp_recvmsg_args, &pid_tgid, &skp, BPF_ANY);
    return 0;
}

SEC("kprobe/tcp_recvmsg")
int BPF_KPROBE_INSTR(WITH(BYPASS), kprobe__tcp_recvmsg__pre_4_1_0) {
    u64 pid_tgid = bpf_get_current_pid_tgid();
    log_debug("kprobe/tcp_recvmsg: pid_tgid: %llu", pid_tgid);
    int flags = (int)PT_REGS_PARM6(ctx);
    if (flags & MSG_PEEK) {
        return 0;
    }

    struct sock *skp = (struct sock*)PT_REGS_PARM2(ctx);
    bpf_map_update_with_telemetry(tcp_recvmsg_args, &pid_tgid, &skp, BPF_ANY);
    return 0;
}

#endif // COMPILE_CORE || COMPILE_PREBUILT

SEC("kretprobe/tcp_recvmsg")
int BPF_KRETPROBE_INSTR(WITH(BYPASS), kretprobe__tcp_recvmsg, int recv) {
    u64 pid_tgid = bpf_get_current_pid_tgid();
    struct sock **skpp = (struct sock**) bpf_map_lookup_elem(&tcp_recvmsg_args, &pid_tgid);
    if (!skpp) {
        return 0;
    }

    struct sock *skp = *skpp;
    bpf_map_delete_elem(&tcp_recvmsg_args, &pid_tgid);
    if (!skp) {
        return 0;
    }

    if (recv < 0) {
        return 0;
    }

    return handle_tcp_recv(pid_tgid, skp, recv);
}

SEC("kprobe/tcp_read_sock")
int BPF_KPROBE_INSTR(WITH(BYPASS), kprobe__tcp_read_sock, struct sock *skp) {
    u64 pid_tgid = bpf_get_current_pid_tgid();
    // we reuse tcp_recvmsg_args here since there is no overlap
    // between the tcp_recvmsg and tcp_read_sock paths
    bpf_map_update_with_telemetry(tcp_recvmsg_args, &pid_tgid, &skp, BPF_ANY);
    return 0;
}

SEC("kretprobe/tcp_read_sock")
int BPF_KRETPROBE_INSTR(WITH(BYPASS), kretprobe__tcp_read_sock, int recv) {
    u64 pid_tgid = bpf_get_current_pid_tgid();
    // we reuse tcp_recvmsg_args here since there is no overlap
    // between the tcp_recvmsg and tcp_read_sock paths
    struct sock **skpp = (struct sock**) bpf_map_lookup_elem(&tcp_recvmsg_args, &pid_tgid);
    if (!skpp) {
        return 0;
    }

    struct sock *skp = *skpp;
    bpf_map_delete_elem(&tcp_recvmsg_args, &pid_tgid);
    if (!skp) {
        return 0;
    }

    if (recv < 0) {
        return 0;
    }

    return handle_tcp_recv(pid_tgid, skp, recv);
}

#endif // __TCP_RECV_H__
