import z from "zod";

export const memberRoleSchema = z.enum(["admin", "user"]);
export type MemberRole = z.infer<typeof memberRoleSchema>;
