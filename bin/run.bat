start /b reg.exe -Config=./config/reg_server.cfg

start /b db.exe -Config=./config/db_server.cfg

start /b serviceUse.exe -Config=./config/service/service_use.cfg

start /b webUseCenter.exe -Config=./config/web/web_center_use.cfg

start /b gateWeb.exe -Config=./config/gate/gate_web.cfg