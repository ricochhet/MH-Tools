#include <koalabox/util.hpp>
#include <koalabox/win_util.hpp>
#include <koalabox/logger.hpp>
#include <koalabox/globals.hpp>

namespace koalabox::util {

    KOALABOX_API(void) error_box(const String& title, const String& message) {
        ::MessageBox(
            nullptr,
            to_wstring(message).c_str(),
            to_wstring(title).c_str(),
            MB_OK | MB_ICONERROR
        );
    }

    [[noreturn]] KOALABOX_API(void) panic(String message) {
        const auto title = fmt::format("[{}] Panic!", globals::get_project_name(false));

        const auto last_error = ::GetLastError();
        if (last_error != 0) {
            message += fmt::format(
                "\n———————— Windows Last Error ————————\nCode: {}\nMessage: {}",
                last_error, win_util::format_message(last_error)
            );
        }

        LOG_CRITICAL("{}", message)

        error_box(title, message);

        exit(static_cast<int>(last_error));
    }

    KOALABOX_API(String) to_string(const WideString& wstr) {
        if (wstr.empty()) {
            return {};
        }

        const auto required_size = WideCharToMultiByte(
            CP_UTF8, 0, wstr.data(), static_cast<int>(wstr.size()), nullptr, 0, nullptr, nullptr
        );

        String string(required_size, 0);
        WideCharToMultiByte(
            CP_UTF8, 0, wstr.data(), static_cast<int>(wstr.size()), string.data(), required_size,
            nullptr, nullptr
        );

        return string;
    }

    KOALABOX_API(WideString) to_wstring(const String& str) {
        if (str.empty()) {
            return {};
        }

        const auto required_size = MultiByteToWideChar(
            CP_UTF8, 0, str.data(), static_cast<int>(str.size()), nullptr, 0
        );

        WideString wstring(required_size, 0);
        MultiByteToWideChar(CP_UTF8, 0, str.data(), static_cast<int>(str.size()), wstring.data(),
            required_size);

        return wstring;
    }

    // Source: https://guidedhacking.com/threads/testing-if-pointer-is-invalid.13222/post-77709
    KOALABOX_API(bool) is_valid_pointer(const void* pointer) {
        const auto mbi_opt = win_util::virtual_query(pointer);

        if (mbi_opt) {
            const auto is_rwe = mbi_opt->Protect & (
                PAGE_READONLY |
                PAGE_READWRITE |
                PAGE_WRITECOPY |
                PAGE_EXECUTE_READ |
                PAGE_EXECUTE_READWRITE |
                PAGE_EXECUTE_WRITECOPY
            );

            const auto is_guarded = mbi_opt->Protect & (
                PAGE_GUARD |
                PAGE_NOACCESS
            );

            return is_rwe && !is_guarded;
        }

        return false;
    }

}
