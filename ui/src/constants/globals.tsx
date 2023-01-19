import { BiCog, BiChip, BiServer } from "react-icons/bi";
import { RiProfileLine } from "react-icons/ri";
import { IconType } from "react-icons/lib";

export function getPage(page: string) {
  return Pages.find((x) => x.page === page);
}

export type Page = {
  page: string;
  icon: IconType;
  verbose: string;
};

export const Pages: Page[] = [
  //{ page: "Home", icon: VscHome, verbose: "Home" },
  { page: "Instance", icon: BiChip, verbose: "Instance" },
  { page: "Profiles", icon: RiProfileLine, verbose: "Profile" },
  { page: "Services", icon: BiServer, verbose: "Service" },
  { page: "Settings", icon: BiCog, verbose: "Setting" },
];

export const customTheme = {
  control: (styles: any) => ({
    ...styles,
    backgroundColor: "#3D3D3D",
    borderColor: "#6B6B6B",
  }),
  input: (styles: any) => ({ ...styles, color: "white" }),
  placeholder: (styles: any) => ({ ...styles, color: "white" }),
  singleValue: (styles: any) => ({
    ...styles,
    color: "white",
  }),
  menu: (styles: any) => ({ ...styles, backgroundColor: "#3D3D3D" }),
  option: (styles: any, { data, isDisabled, isFocused, isSelected }: any) => ({
    ...styles,
    backgroundColor: isSelected ? "#525252" : isFocused ? "#00A3FF" : "#3D3D3D",
  }),
};

export type InteractionOption = {
  value: string;
  label: string;
};

export const InteractionOptions: InteractionOption[] = [
  { value: "high", label: "High" },
  { value: "low", label: "Low" },
];

export type NetworkOption = {
  value: string;
  label: string;
};

export const NetworkOptions: NetworkOption[] = [
  { value: "tcp", label: "TCP" },
  { value: "udp", label: "UDP" },
];
