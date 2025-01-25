import { get, post } from "@/utils/networkWrapper.js";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const useStatus = () => {
  return useQuery({
    queryKey: ["status"],
    queryFn: async () => {
      const response = await get("http://localhost:4000/api/scraper/status");
      console.log(response);
      return {
        isRunning: response.isRunning,
        interval: response.interval,
      };
    },
  });
};

export function usePostStartScraper() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["startScraper"],
    mutationFn: async () => {
      return await post("http://localhost:4000/api/scraper/start");
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
      return await post("http://localhost:4000/api/scraper/stop");
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
      return await post("http://localhost:4000/api/scraper/force-start");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["status"] });
    },
  });
}

export function usePostRestartScraper() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["restartScraper"],
    mutationFn: async () => {
      return await post("http://localhost:4000/api/scraper/restart");
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
      return await post("http://localhost:4000/api/scraper/change-interval", {
        interval,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["status"] });
    },
  });
}
