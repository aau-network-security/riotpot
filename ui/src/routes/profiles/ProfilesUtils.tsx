import { CreateButton, UtilsBar } from "../../components/utils/Utils";
import { getPage } from "../../constants/globals";
import ProfileForm from "./ProfileForm";

export const ProfilesUtils = () => {
  // Utils buttons
  const page = getPage("Profiles");
  const content = <ProfileForm create={false} />;
  const createbtn = (
    <CreateButton
      title="Profile"
      icon={page && page.icon}
      content={content}
    ></CreateButton>
  );
  const buttons = [createbtn];
  return <UtilsBar buttons={buttons} />;
};
