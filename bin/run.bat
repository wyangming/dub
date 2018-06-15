start /b reg.exe -Config=./config/reg_server.cfg

start /b db.exe -Config=./config/db_server.cfg

start /b serviceUse.exe -Config=./config/service/service_use.cfg

start /b manLobby.exe -Config=./config/web/web_center_man_lobby.cfg

start /b manAgent.exe -Config=./config/web/web_center_man_agent.cfg

start /b gateWeb.exe -Config=./config/gate/gate_web.cfg