#ifndef __SHARED_LIBRARIES_PROBES_H
#define __SHARED_LIBRARIES_PROBES_H

#include "bpf_telemetry.h"
#include "bpf_bypass.h"

#include "pid_tgid.h"
#include "shared-libraries/types.h"

static __always_inline void fill_path_safe(lib_path_t *path, const char *path_argument) {
#pragma unroll
    for (int i = 0; i < LIB_PATH_MAX_SIZE; i++) {
        bpf_probe_read_user(&path->buf[i], 1, &path_argument[i]);
        if (path->buf[i] == 0) {
            path->len = i;
            break;
        }
    }
}

static __always_inline void do_sys_open_helper_enter(const char *filename) {
    lib_path_t path = { 0 };
    if (bpf_probe_read_user_with_telemetry(path.buf, sizeof(path.buf), filename) >= 0) {
// Find the null character and clean up the garbage following it
#pragma unroll
        for (int i = 0; i < LIB_PATH_MAX_SIZE; i++) {
            if (path.buf[i] == 0) {
                path.len = i;
                break;
            }
        }
    } else {
        fill_path_safe(&path, filename);
    }

    // Bail out if the path size is larger than our buffer
    if (!path.len) {
        return;
    }

    u64 pid_tgid = bpf_get_current_pid_tgid();
    path.pid = GET_USER_MODE_PID(pid_tgid);
    bpf_map_update_with_telemetry(open_at_args, &pid_tgid, &path, BPF_ANY);
    return;
}

static __always_inline void do_sys_open_helper_exit(exit_sys_ctx *args) {
    u64 pid_tgid = bpf_get_current_pid_tgid();

    // If file couldn't be opened, bail out
    if (args->ret < 0) {
        goto cleanup;
    }

    lib_path_t *path = bpf_map_lookup_elem(&open_at_args, &pid_tgid);
    if (path == NULL) {
        return;
    }

    // Check the last 9 characters of the following libraries to ensure the file is a relevant `.so`.
    // Libraries:
    //    libssl.so -> libssl.so
    // libcrypto.so -> crypto.so
    // libgnutls.so -> gnutls.so
    //
    // The matching is done in 2 stages here, first we look if the filename finished by ".so" 6 chars forward
    // this will give us the index (where the loop) for the 2nd stage
    // 2nd stage will try to match the remaining
    // it's done this way to avoid unroll code generation complexity and some verifier don't allow that
    bool is_shared_library = false;
#define match3chars(_base, _a, _b, _c) (path->buf[_base + i] == _a && path->buf[_base + i + 1] == _b && path->buf[_base + i + 2] == _c)
#define match6chars(_base, _a, _b, _c, _d, _e, _f) (match3chars(_base, _a, _b, _c) && match3chars(_base + 3, _d, _e, _f))
    int i = 0;
#pragma unroll
    for (i = 0; i < LIB_PATH_MAX_SIZE - (LIB_SO_SUFFIX_SIZE); i++) {
        if (match3chars(6, '.', 's', 'o')) {
            is_shared_library = true;
            break;
        }
    }
    if (!is_shared_library) {
        goto cleanup;
    }

    u64 crypto_libset_enabled = 0;
    LOAD_CONSTANT("crypto_libset_enabled", crypto_libset_enabled);

    if (crypto_libset_enabled && (match6chars(0, 'l', 'i', 'b', 's', 's', 'l') || match6chars(0, 'c', 'r', 'y', 'p', 't', 'o') || match6chars(0, 'g', 'n', 'u', 't', 'l', 's'))) {
        bpf_perf_event_output((void *)args, &crypto_shared_libraries, BPF_F_CURRENT_CPU, path, sizeof(lib_path_t));
        goto cleanup;
    }

    u64 gpu_libset_enabled = 0;
    LOAD_CONSTANT("gpu_libset_enabled", gpu_libset_enabled);

    if (gpu_libset_enabled && (match6chars(0, 'c', 'u', 'd', 'a', 'r', 't'))) {
        bpf_perf_event_output((void *)args, &gpu_shared_libraries, BPF_F_CURRENT_CPU, path, sizeof(lib_path_t));
    }

cleanup:
    bpf_map_delete_elem(&open_at_args, &pid_tgid);
    return;
}

// This definition is the same for all architectures.
#ifndef O_WRONLY
#define O_WRONLY        00000001
#endif

static __always_inline int should_ignore_flags(int flags)
{
    return flags & O_WRONLY;
}

SEC("tracepoint/syscalls/sys_enter_open")
int tracepoint__syscalls__sys_enter_open(enter_sys_open_ctx *args) {
    CHECK_BPF_PROGRAM_BYPASSED()

    if (should_ignore_flags(args->flags)) {
        return 0;
    }

    do_sys_open_helper_enter(args->filename);
    return 0;
}

SEC("tracepoint/syscalls/sys_exit_open")
int tracepoint__syscalls__sys_exit_open(exit_sys_ctx *args) {
    CHECK_BPF_PROGRAM_BYPASSED()
    do_sys_open_helper_exit(args);
    return 0;
}

SEC("tracepoint/syscalls/sys_enter_openat")
int tracepoint__syscalls__sys_enter_openat(enter_sys_openat_ctx *args) {
    CHECK_BPF_PROGRAM_BYPASSED()

    if (should_ignore_flags(args->flags)) {
        return 0;
    }

    do_sys_open_helper_enter(args->filename);
    return 0;
}

SEC("tracepoint/syscalls/sys_exit_openat")
int tracepoint__syscalls__sys_exit_openat(exit_sys_ctx *args) {
    CHECK_BPF_PROGRAM_BYPASSED()
    do_sys_open_helper_exit(args);
    return 0;
}

SEC("tracepoint/syscalls/sys_enter_openat2")
int tracepoint__syscalls__sys_enter_openat2(enter_sys_openat2_ctx *args) {
    CHECK_BPF_PROGRAM_BYPASSED()
    // Unlike the other variants, openat2(2) has the flags embedded inside the
    // how argument; we don't bother trying to accessing it for now.
    do_sys_open_helper_enter(args->filename);
    return 0;
}

SEC("tracepoint/syscalls/sys_exit_openat2")
int tracepoint__syscalls__sys_exit_openat2(exit_sys_ctx *args) {
    CHECK_BPF_PROGRAM_BYPASSED()
    do_sys_open_helper_exit(args);
    return 0;
}

#endif
