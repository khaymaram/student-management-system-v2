import { useQuery } from "@tanstack/react-query";

import { get } from "../lib/axios";
import type { Major } from "../types";

export function useMajors() {
  return useQuery<Major[]>({
    queryKey: ["majors"],
    queryFn: () => get<Major[]>("/majors"),
    staleTime: 5 * 60 * 1000,
  });
}
