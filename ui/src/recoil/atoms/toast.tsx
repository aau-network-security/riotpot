import { AiOutlineInfoCircle } from "react-icons/ai";
import { atom } from "recoil";

export type ToastVariant = {
  color: string;
  icon: JSX.Element | undefined;
  name: string;
};

export const InfoToastVariant = {
  color: "light",
  icon: <AiOutlineInfoCircle />,
  name: "Info",
} as ToastVariant;

export const ErrorToastVariant = {
  color: "danger",
  icon: <AiOutlineInfoCircle />,
  name: "Error",
} as ToastVariant;

export type ToastState = {
  variant: ToastVariant;
  title: string;
  message: string;
  show: boolean;
};

export const InfoToast = {
  variant: InfoToastVariant,
  title: "",
  message: "",
  show: false,
} as ToastState;

export const toast = atom<ToastState>({
  key: "toast",
  default: InfoToast,
});
