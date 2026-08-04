
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { get, put } from "../lib/axios";
import type { Finance } from "../types";

export function useFinances() {
    return useQuery<Finance[]>({
        queryKey: ["finances"],
        queryFn: async () => await get<Finance[]>("/finances"),
    });
}

export function useFinance(id: number | null) {
    return useQuery<Finance>({
        queryKey: ["finance", id],
        queryFn: () => get<Finance>(`/finances/${id}`),
        enabled: id !== null,
    });
}

export function useUpdateFinance() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({
            studentId,
            input,
        }: {
            studentId: number;
            input: {
                scholarship: number;
                paid: number;
                isInState: boolean;
            };
        }) => {
            return await put<string>(`/finances/${studentId}`, input);
        },

        onSuccess: (_, variables) => {
            // Refresh the finance list
            queryClient.invalidateQueries({
                queryKey: ["finances"],
            });

            // Refresh the individual student's finance details
            queryClient.invalidateQueries({
                queryKey: ["finance", variables.studentId],
            });
        },
    });
}
