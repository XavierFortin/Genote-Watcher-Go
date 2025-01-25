import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useMemo } from "react";
import { useState } from "react";

export const useLogs = () => {
  return useQuery({
    queryKey: ["logs"],
    queryFn: async () => {
      const response = await axios.get("/api/logs");
      return (await response.data()).split("\n");
    },
  });
};

export const useWebsocketLogger = () => {
  const memoizedWs = useMemo(() => new WebSocket("/ws"), []);
  const [logs, setLogs] = useState([]);

  memoizedWs.onmessage = (event) => {
    setLogs((prevLogs) => [...prevLogs, event.data]);
  };

  return logs;
};
