#include <smoke_api/config.hpp>
#include <core/paths.hpp>
#include <koalabox/util.hpp>
#include <koalabox/io.hpp>
#include <koalabox/logger.hpp>

namespace smoke_api::config {
    Config instance; // NOLINT(cert-err58-cpp)

    // TODO: Refactor to Koalabox
    void init_config() {
        const auto path = paths::get_config_path();

        if (exists(path)) {
            try {
                const auto config_str = koalabox::io::read_file(path);

                instance = Json::parse(config_str).get<Config>();

                LOG_DEBUG("Parsed config:\n{}", Json(instance).dump(2))
            } catch (const Exception& e) {
                const auto message = fmt::format("Error parsing config file: {}", e.what());
                koalabox::util::error_box("SmokeAPI Error", message);
            }
        }
    }

    Vector<DLC> get_extra_dlcs(AppId_t app_id) {
        return DLC::get_dlcs_from_apps(instance.extra_dlcs, app_id);
    }

    bool is_dlc_unlocked(AppId_t app_id, AppId_t dlc_id, const Function<bool()>& original_function) {
        auto status = instance.default_app_status;

        const auto app_id_str = std::to_string(app_id);
        if (instance.override_app_status.contains(app_id_str)) {
            status = instance.override_app_status[app_id_str];
        }

        const auto dlc_id_str = std::to_string(dlc_id);
        if (instance.override_dlc_status.contains(dlc_id_str)) {
            status = instance.override_dlc_status[dlc_id_str];
        }

        bool is_unlocked;
        switch (status) {
            case AppStatus::UNLOCKED:
                is_unlocked = true;
                break;
            case AppStatus::LOCKED:
                is_unlocked = false;
                break;
            case AppStatus::ORIGINAL:
            case AppStatus::UNDEFINED:
                is_unlocked = original_function();
                break;
        }

        LOG_TRACE(
            "App ID: {}, DLC ID: {}, Status: {}, Original: {}, Unlocked: {}",
            app_id_str, dlc_id_str, Json(status).dump(), original_function(), is_unlocked
        )

        return is_unlocked;
    }

    DLL_EXPORT(void) ReloadConfig() {
        LOG_INFO("Reloading config")

        init_config();
    }
}
