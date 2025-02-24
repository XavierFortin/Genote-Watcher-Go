import { useLogs, useWebsocketLogger } from "@/services/hooks/logger.js";
import "./LogViewer.css";

export function LogViewer() {
  var { data: currentLogs, isPending } = useLogs();
  var newLogs = useWebsocketLogger();
  return (
    <div className="log-viewer">
      {!isPending &&
        currentLogs.map((log, index) => (
          <div key={index} className="log-line">
            {log}
          </div>
        ))}
      {newLogs.map((log, index) => (
        <div key={index} className="log-line">
          {log}
        </div>
      ))}
    </div>
  );
}
