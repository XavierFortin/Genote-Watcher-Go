import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export const useStatus = () => {
  return useQuery({
    queryKey: ["status"],
    queryFn: async () => {
      const response = await axios.get("/api/scraper/status");
      return {
        isRunning: response.data.isRunning,
        interval: response.data.interval,
      };
    },
  });
};

export function usePostStartScraper() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["startScraper"],
    mutationFn: async () => {
      return await axios.post("/api/scraper/start");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["status"] });
    },
  });
}

export function usePostStopScraper() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["stopScraper"],
    mutationFn: async () => {
      return await axios.post("/api/scraper/stop");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["status"] });
    },
  });
}

export function usePostForceStartOnceScraper() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["forceStartOnceScraper"],
    mutationFn: async () => {
      return await axios.post("/api/scraper/force-start");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["status"] });
    },
  });
}

export function usePostChangeInterval() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["changeInterval"],
    mutationFn: async (interval) => {
      return await axios.post("/api/scraper/change-interval", {
        interval,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["status"] });
    },
  });
}
