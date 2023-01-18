import { Stack } from "react-bootstrap";
import { SettingsAPIAddress } from "./SettingsBackend";
import { SettingsState } from "./SettingsState";

export const Settings = () => {
  return (
    <Stack gap={5}>
      <SettingsAPIAddress />
      <SettingsState />
    </Stack>
  );
};
