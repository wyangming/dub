"">./log/webCenterManLobby.log
"">./log/webCenterManAgent.log

"">./log/serviceUse.log

"">./log/db.log

"">./log/reg.log

"">./log/gateWeb.log

::clear app
del /F /S /Q db.exe
del /F /S /Q gateWeb.exe
del /F /S /Q manLobby.exe
del /F /S /Q manAgent.exe
del /F /S /Q reg.exe
del /F /S /Q serviceUse.exe