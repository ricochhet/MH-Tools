#include <store_mode/steamclient/steamclient.hpp>
#include <steam_impl/steam_apps.hpp>

VIRTUAL(bool) IClientAppManager_IsAppDlcInstalled(PARAMS(AppId_t app_id, AppId_t dlc_id)) {
    try {
        return steam_apps::IsDlcUnlocked(
            __func__, app_id, dlc_id, [&]() {
                GET_ORIGINAL_HOOKED_FUNCTION(IClientAppManager_IsAppDlcInstalled)

                return IClientAppManager_IsAppDlcInstalled_o(ARGS(app_id, dlc_id));
            }
        );
    } catch (const Exception& e) {
        LOG_ERROR("{} -> Error: {}", __func__, e.what())
        return false;
    }
}
