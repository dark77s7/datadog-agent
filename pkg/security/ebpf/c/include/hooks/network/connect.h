#ifndef _HOOKS_CONNECT_H_
#define _HOOKS_CONNECT_H_

#include "constants/offsets/netns.h"
#include "constants/syscall_macro.h"
#include "helpers/discarders.h"

HOOK_SYSCALL_ENTRY3(connect, int, socket, struct sockaddr *, addr, unsigned int, addr_len) {
    if (!addr) {
        return 0;
    }

    struct syscall_cache_t syscall = {
        .type = EVENT_CONNECT,
    };
    cache_syscall(&syscall);
    return 0;
}

int __attribute__((always_inline)) sys_connect_ret(void *ctx, int retval) {
    struct syscall_cache_t *syscall = pop_syscall(EVENT_CONNECT);
    if (!syscall) {
        return 0;
    }

    if (IS_UNHANDLED_ERROR(retval)) {
        return 0;
    }

    /* pre-fill the event */
    struct connect_event_t event = {
        .syscall.retval = retval,
        .addr[0] = syscall->connect.addr[0],
        .addr[1] = syscall->connect.addr[1],
        .family = syscall->connect.family,
        .port = syscall->connect.port,
        .protocol = syscall->connect.protocol,
    };

    struct proc_cache_t *entry = fill_process_context(&event.process);
    fill_container_context(entry, &event.container);
    fill_span_context(&event.span);

    // Check if we should sample this event for activity dumps
    struct activity_dump_config *config = lookup_or_delete_traced_pid(event.process.pid, bpf_ktime_get_ns(), NULL);
    if (config) {
        if (mask_has_event(config->event_mask, EVENT_CONNECT)) {
            event.event.flags |= EVENT_FLAGS_ACTIVITY_DUMP_SAMPLE;
        }
    }

    send_event(ctx, EVENT_CONNECT, event);
    return 0;
}

HOOK_SYSCALL_EXIT(connect) {
    int retval = SYSCALL_PARMRET(ctx);
    return sys_connect_ret(ctx, retval);
}

HOOK_ENTRY("security_socket_connect")
int hook_security_socket_connect(ctx_t *ctx) {
    struct socket *sk = (struct socket *)CTX_PARM1(ctx);
    struct sockaddr *address = (struct sockaddr *)CTX_PARM2(ctx);
    struct pid_route_t key = {};
    u16 family = 0;
    u16 protocol = 0;
    
    u64 socket_sock_offset;
    u64 sk_protocol_offset;
    u64 sk_protocol_size;


    LOAD_CONSTANT("socket_sock_offset", socket_sock_offset);
    LOAD_CONSTANT("sock_sk_protocol_offset", sk_protocol_offset);
    LOAD_CONSTANT("sk_protocol_size", sk_protocol_size);

    __bpf_printk("-------------------- sk_protocol_offset: %d", socket_sock_offset);
    __bpf_printk("-------------------- sk_protocol_offset: %d", sk_protocol_offset);
    __bpf_printk("-------------------- sk_protocol_size: %d", sk_protocol_size);


    // Extract IP and port from the sockaddr structure
    bpf_probe_read(&family, sizeof(family), &address->sa_family);

    if (family == AF_INET) {
        struct sockaddr_in *addr_in = (struct sockaddr_in *)address;
        bpf_probe_read(&key.port, sizeof(addr_in->sin_port), &addr_in->sin_port);
        bpf_probe_read(&key.addr, sizeof(addr_in->sin_addr.s_addr), &addr_in->sin_addr.s_addr);
    } else if (family == AF_INET6) {
        struct sockaddr_in6 *addr_in6 = (struct sockaddr_in6 *)address;
        bpf_probe_read(&key.port, sizeof(addr_in6->sin6_port), &addr_in6->sin6_port);
        bpf_probe_read(&key.addr, sizeof(u64) * 2, (char *)addr_in6 + offsetof(struct sockaddr_in6, sin6_addr));
    }

    struct sock *sk_sock = NULL;
    bpf_probe_read(&sk_sock, sizeof(sk_sock),(void *) sk + socket_sock_offset);
    bpf_probe_read(&protocol, sk_protocol_size, (void *) sk_sock + sk_protocol_offset);
    // bpf_probe_read(&protocol, sizeof(protocol), &sk_sock->sk_protocol);

    // fill syscall_cache if necessary
    struct syscall_cache_t *syscall = peek_syscall(EVENT_CONNECT);
    if (syscall) {
        syscall->connect.addr[0] = key.addr[0];
        syscall->connect.addr[1] = key.addr[1];
        syscall->connect.port = key.port;
        syscall->connect.family = family;
        syscall->connect.protocol = protocol;
    }

    // Only handle AF_INET and AF_INET6
    if (family != AF_INET && family != AF_INET6) {
        return 0;
    }

    return 0;
}

SEC("tracepoint/handle_sys_connect_exit")
int tracepoint_handle_sys_connect_exit(struct tracepoint_raw_syscalls_sys_exit_t *args) {
    return sys_connect_ret(args, args->ret);
}

#endif /* _CONNECT_H_ */
