import {
  Address,
  DeleteDropdownItem,
  InteractionBadge,
  NetworkBadge,
  OptionsDropdown,
} from "../../components/utils/Common";
import {
  getPage,
  InteractionOption,
  InteractionOptions,
  NetworkOption,
  NetworkOptions,
} from "../../constants/globals";
import Table from "../../components/table/Table";

export type ServicesProps = {
  service: string;
  interaction: InteractionOption;
  network: NetworkOption;
  host: string;
  port: Number;
};

export const ServiceRow = ({
  service,
  interaction,
  network,
  host,
  port,
}: ServicesProps) => {
  const intBadge = <InteractionBadge {...interaction} />;
  const netBadge = <NetworkBadge {...network} />;
  const address = <Address port={port} host={host} />;

  const page = getPage("Services");
  const note =
    "This service will also be deleted from other instances using this service.";
  const options = (
    <OptionsDropdown>
      {page && <DeleteDropdownItem page={page} note={note} name={service} />}
    </OptionsDropdown>
  );

  return [service, intBadge, netBadge, address, options];
};

export const ServicesTable = () => {
  const sprops: ServicesProps = {
    service: "CoAP",
    interaction: InteractionOptions[1],
    network: NetworkOptions[1],
    host: "localhost",
    port: 1111,
  };

  // Mock Data
  const data = {
    headers: ["Service", "Interaction", "Network", "Address", ""],
    rows: [ServiceRow(sprops), ServiceRow(sprops)],
  };

  return <Table data={data} />;
};
