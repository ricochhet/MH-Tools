#include <koalabox/crypto.hpp>
#include <koalabox/logger.hpp>
#include <koalabox/win_util.hpp>
#include <koalabox/util.hpp>

#include <wincrypt.h>

namespace koalabox::crypto {

    Vector<uint8_t> decode_hex_string(const String& hex_str) {
        if (hex_str.length() < 2) {
            return {};
        }

        Vector<uint8_t> buffer(hex_str.size() / 2);

        std::stringstream ss;
        ss << std::hex << hex_str;

        for (size_t i = 0; i < hex_str.length(); i++) {
            ss >> buffer[i];
        }

        return buffer;
    }

    // Source: https://learn.microsoft.com/en-us/windows/win32/seccrypto/example-c-program--creating-an-md-5-hash-from-file-content
    String calculate_md5(const Path& file_path) {
        String result_buffer(32, '\0');
        const auto buffer_size = 1024 * 1024; // 1 Mb
        const auto md5_len = 16;

        BOOL result = FALSE;
        HCRYPTPROV hProv = 0;
        HCRYPTHASH hHash = 0;
        HANDLE hFile = NULL;

        hFile = CreateFile(
            file_path.wstring().c_str(),
            GENERIC_READ,
            FILE_SHARE_READ,
            NULL,
            OPEN_EXISTING,
            FILE_FLAG_SEQUENTIAL_SCAN,
            NULL
        );

        if (INVALID_HANDLE_VALUE == hFile) {
            LOG_ERROR(
                "Error opening file {}. Error: {}",
                file_path.string(), win_util::get_last_error()
            )

            return "";
        }

        // Get handle to the crypto provider
        if (!CryptAcquireContext(
            &hProv,
            NULL,
            NULL,
            PROV_RSA_FULL,
            CRYPT_VERIFYCONTEXT
        )) {
            LOG_ERROR("CryptAcquireContext error. Error: {}", win_util::get_last_error())

            CloseHandle(hFile);
            return "";
        }

        if (!CryptCreateHash(hProv, CALG_MD5, 0, 0, &hHash)) {
            LOG_ERROR("CryptAcquireContext error. Error: {}", win_util::get_last_error())

            CloseHandle(hFile);
            CryptReleaseContext(hProv, 0);
            return "";
        }

        BYTE rgb_file[buffer_size];
        DWORD bytes_read = 0;
        while ((result = ReadFile(
            hFile,
            rgb_file,
            buffer_size,
            &bytes_read,
            NULL
        ))) {
            if (not bytes_read) {
                break;
            }

            if (!CryptHashData(hHash, rgb_file, bytes_read, 0)) {
                LOG_ERROR("CryptHashData error. Error: {}", win_util::get_last_error())

                CryptReleaseContext(hProv, 0);
                CryptDestroyHash(hHash);
                CloseHandle(hFile);
                return "";
            }
        }

        if (!result) {
            LOG_ERROR("ReadFile error. Error: {}", win_util::get_last_error())

            CryptReleaseContext(hProv, 0);
            CryptDestroyHash(hHash);
            CloseHandle(hFile);
            return "";
        }

        BYTE rgb_hash[md5_len];
        CHAR rgb_digits[] = "0123456789abcdef";
        DWORD hash_length = md5_len;
        if (CryptGetHashParam(hHash, HP_HASHVAL, rgb_hash, &hash_length, 0)) {
            for (DWORD i = 0; i < hash_length; i++) {
                result_buffer[i * 2] = rgb_digits[rgb_hash[i] >> 4];
                result_buffer[i * 2 + 1] = rgb_digits[rgb_hash[i] & 0xf];
            }
        } else {
            LOG_ERROR("ReadFile CryptGetHashParam. Error: {}", win_util::get_last_error())
        }

        CryptDestroyHash(hHash);
        CryptReleaseContext(hProv, 0);
        CloseHandle(hFile);

        return result_buffer;
    }

}