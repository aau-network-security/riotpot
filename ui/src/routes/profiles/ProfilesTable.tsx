import Table from "../../components/table/Table";
import {
  DeleteDropdownItem,
  OptionsDropdown,
} from "../../components/utils/Common";
import {
  getPage,
  InteractionOption,
  InteractionOptions,
  NetworkOption,
  NetworkOptions,
} from "../../constants/globals";
import { ProfileRowInfoPop } from "../instances/InstancesUtils";

export interface ProfileService {
  name: string;
  interaction: InteractionOption;
  network: NetworkOption;
}

export interface Profile {
  name: string;
  services: ProfileService[];
}

const ProfileRowOptions = ({ name }: { name: string }) => {
  const page = getPage("Profiles");
  const note =
    "This profile will be unregistered. However, instances using this profile will not stop working";
  return (
    <OptionsDropdown>
      {page && <DeleteDropdownItem page={page} note={note} name={name} />}
    </OptionsDropdown>
  );
};

const ProfileRow = ({ name, services }: Profile) => {
  return [
    name,
    <ProfileRowInfoPop services={services} />,
    <ProfileRowOptions name={name} />,
  ];
};

export const ProfilesTable = () => {
  const profile = {
    name: "Wi-Fi Printer",
    services: [
      {
        name: "CoAP",
        interaction: InteractionOptions[0],
        network: NetworkOptions[0],
      },
    ],
  };

  const rows = [ProfileRow(profile), ProfileRow(profile)];

  const data = {
    headers: [`${rows.length} Profiles`, "", ""],
    rows: rows,
  };

  return <Table data={data} />;
};
