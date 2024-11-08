if not fs.exists("json.lua") then
  shell.run("wget","https://raw.githubusercontent.com/rxi/json.lua/master/json.lua")
end

local json = require "json"

if not fs.exists("id") then
    res = http.get("https://www.uuidtools.com/api/generate/v4")
    local content = res.readAll()
    local id = json.decode(content)[1]
    local file = fs.open("id","w")
    file.write(id)
    file.close()
end


function executeCode(codeString)
  local func,err=loadstring(codeString)
  local success=false
  if func then
      success,err = pcall(func)
  end

  return success,err
end

function websocketloop()
  local file = fs.open("id", "r")
  local id = file.readAll()
  file.close()

  local ws, err = http.websocket("wss://cc-remote-control.fly.dev/ws", {["User-Agent"] = id})
  if err then
    print(err)
    os.sleep(30)
    return
  end
  print("Connected as ", id)
  while true do
    print("waiting..")
    local output = ws.receive()

    print(output)
    local obj = json.decode(output)
    if obj['Action'] == "eval" then
        local ok = false
        local err = nil
        ok, err = executeCode(obj['Code'])
        ws.send(err)
    elseif obj['Action'] == "kill" then
        ws.send("bye bye")
        return "kill"
    elseif obj['Action'] == "shell" then
        shell.run(obj['Code'])
        ws.send("No output currently, WIP")
    end
  end
end


while true do
    local status, res = pcall(websocketloop)
    print(res)
    if res == "kill" then
        break
    end
end
