import { useRecoilState, useRecoilValue, useSetRecoilState } from "recoil";
import Table from "../../components/table/Table";
import {
  DeleteDropdownItem,
  EditDropdownItem,
  OptionsDropdown,
} from "../../components/utils/Common";
import { getPage } from "../../constants/globals";
import {
  profiles,
  Profile,
  removeProfile,
  profileFormFieldErrors,
  profilesFilter,
} from "../../recoil/atoms/profiles";

import { ProfileRowInfoPop } from "../instances/InstanceUtils";
import { SimpleForm } from "../../components/forms/Form";
import { ProfileFormFields } from "./ProfileForm";
import { useLocation } from "react-router-dom";
import { CgDetailsLess } from "react-icons/cg";
import { Dropdown } from "react-bootstrap";

const EditProfile = ({ profile }: { profile: Profile }) => {
  // Get the page to set the icon
  const page = getPage("Profiles");

  // Get the current values of the profile
  const setter = useSetRecoilState(profilesFilter(profile.id));

  // Create a setter for the submit
  const onSubmit = (newElement: any) => {
    setter(newElement);
  };

  // Create the form with the default values as they currently are
  const content = (
    <SimpleForm
      create={false}
      defaultValues={profile}
      errors={profileFormFieldErrors}
      onSubmit={onSubmit}
      page="Profiles"
      fields={ProfileFormFields}
    />
  );

  // Get the form with the update tag
  return (
    <EditDropdownItem form={content} icon={page?.icon} title={"profile"} />
  );
};

const ViewProfile = ({ profile }: { profile: Profile }) => {
  const location = useLocation();
  const link = `${location.pathname}/${profile.id}`;

  return (
    <Dropdown.Item href={link}>
      <CgDetailsLess />
      Details
    </Dropdown.Item>
  );
};

const ProfileRowOptions = ({ profile }: { profile: Profile }) => {
  // Profiles
  const [profilesList, setProfilesList] = useRecoilState(profiles);
  // Delete option
  const deleteProfile = () =>
    removeProfile(profile.id, profilesList, setProfilesList);

  const page = getPage("Profiles");
  const note =
    "This profile will be unregistered. However, instances using this profile will not stop working";
  return (
    <OptionsDropdown>
      <ViewProfile profile={profile} />
      <EditProfile profile={profile} />
      {page && (
        <DeleteDropdownItem
          page={page}
          note={note}
          name={profile.name}
          onClick={deleteProfile}
        />
      )}
    </OptionsDropdown>
  );
};

const profileRow = (profile: Profile) => {
  return [
    profile.name,
    <ProfileRowInfoPop key={1} services={profile.services} />,
    <ProfileRowOptions key={2} profile={profile} />,
  ];
};

export const ProfilesTable = () => {
  const profilesList = useRecoilValue(profiles);
  const rows = profilesList.map((profile: Profile) => profileRow(profile));

  const data = {
    headers: [`${rows.length} Profiles`, "", ""],
    rows: rows,
  };

  return <Table data={data} />;
};
