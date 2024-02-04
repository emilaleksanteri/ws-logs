import { useState } from "react"

type Log = {
  domain: string
  accessed_at: string
  is_my_domain: boolean
}

function LogsConnect() {
  const ws = new WebSocket("ws://localhost:8080/logs/channel")
  console.log("connected")
  return (
    <LogsListener ws={ws} />
  )
}

function LogsListener(props: { ws: WebSocket }) {
  const { ws } = props

  const [logs, setLogs] = useState<Log[]>([{
    domain: "example.com",
    accessed_at: new Date().toISOString(),
    is_my_domain: true
  }])

  ws.addEventListener("message", (event) => {
    const data: Log = JSON.parse(event.data)
    setLogs([data, ...logs])
  })

  return (
    <Logs logs={logs} />
  )
}

function Logs(props: { logs: Log[] }) {
  const { logs } = props
  return (
    <ul className="bg-white p-4 w-[60%] flex flex-col gap-2 h-[400px] overflow-y-scroll rounded-md">
      {logs.map((log, i) => (
        <li className="flex gap-4 items-center justify-between font-medium pb-2 border-b-2" key={i}>
          <p>{log.domain}</p>
          <p>{log.accessed_at}</p>
          <p className={`capitalize ${log.is_my_domain ? "text-green-400" : "text-rose-500"}`}>
            {log.is_my_domain ? "success" : "wrong domain"}
          </p>
        </li>
      ))}
    </ul>
  )
}

function App() {
  const [showLogs, setShowLogs] = useState(false)
  return (
    <div className="w-screen flex items-center justify-center">
      <div className="w-[70%] py-10 flex flex-col gap-6">
        <h1 className="text-4xl font-bold">Live Logs</h1>
        <div className="w-full bg-zinc-200 py-6 px-4 rounded-md">
          <button
            className="bg-white px-4 py-2 rounded-md font-semibold"
            onClick={() => setShowLogs(!showLogs)}>
            Toggle Show Logs
          </button>
          {
            showLogs
              ? <div className="py-6">
                <p className="font-thin italic pb-4">Listening for site entries...</p>
                <LogsConnect />
              </div>
              : null
          }
        </div>
      </div>
    </div>
  )
}

export default App
