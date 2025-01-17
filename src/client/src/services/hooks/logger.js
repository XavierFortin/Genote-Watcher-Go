import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
import { useState } from "react";

export const useLogs = () => {
  return useQuery({
    queryKey: ["logs"],
    queryFn: async () => {
      const response = await fetch("http://localhost:4000/api/logs");
      return (await response.text()).split("\n");
    },
  });
};

export const useWebsocketLogger = () => {
  const memoizedWs = useMemo(() => new WebSocket("ws://localhost:4000/ws"), []);
  const [logs, setLogs] = useState([]);

  memoizedWs.onmessage = (event) => {
    setLogs((prevLogs) => [...prevLogs, event.data]);
  };

  return logs;
};
