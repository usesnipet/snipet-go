import { http } from "@/lib/http";

import { updateProfilePictureSchema, userSchema } from "./schemas";

import type { UpdateProfilePicture, User } from "./schemas";
import type { ServiceGetOptions, ServicePutOptions } from "@/lib/services";

const USERS_URL = "/api/users";

const me = async (opts: ServiceGetOptions<User> = {}): Promise<User> => {
  return http.get({
    url: `${USERS_URL}/me`,
    schemas: { response: userSchema },
    ...opts,
  });
};

const updatePicture = async (
  body: UpdateProfilePicture,
  opts: ServicePutOptions<UpdateProfilePicture, void> = {},
): Promise<void> => {
  return http.put({
    url: `${USERS_URL}/me/picture`,
    body,
    schemas: {
      body: updateProfilePictureSchema,
    },
    ...opts,
  });
};

export const userService = {
  me,
  updatePicture,
};
