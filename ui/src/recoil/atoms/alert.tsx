import { atom } from "recoil";

export type alertType = {
  type: string;
  color: string;
};

export const alertTypes: alertType[] = [
  { type: "error", color: "" },
  { type: "success", color: "" },
  { type: "warning", color: "" },
];

export const alert = atom({
  key: "alet",
  default: {
    message: "",
    type: "error",
  },
});
