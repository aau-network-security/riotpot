import {
  Address,
  DeleteDropdownItem,
  InteractionBadge,
  NetworkBadge,
  OptionsDropdown,
} from "../../components/utils/Common";
import { getPage } from "../../constants/globals";
import Table from "../../components/table/Table";
import { services, Service } from "../../recoil/atoms/services";
import { useRecoilState, useRecoilValue } from "recoil";
import { removeFromList } from "../../recoil/atoms/common";

const ServiceRowOptions = ({ service }: { service: Service }) => {
  // Services
  const [serviceGetter, serviceSetter] = useRecoilState(services);
  // Delete option
  const deleteService = () =>
    removeFromList(service.id, serviceGetter, serviceSetter);

  const page = getPage("Services");
  const note =
    "This service will also be deleted from other instances using this service.";

  return (
    <OptionsDropdown>
      {page && (
        <DeleteDropdownItem
          page={page}
          note={note}
          name={service.name}
          onClick={deleteService}
        />
      )}
    </OptionsDropdown>
  );
};

const serviceRow = (service: Service) => {
  return [
    service.name,
    <InteractionBadge {...service.interaction} />,
    <NetworkBadge {...service.network} />,
    <Address port={service.port} host={service.host} />,
    <ServiceRowOptions service={service} />,
  ];
};

export const ServicesTable = () => {
  const servicesList = useRecoilValue(services);
  const rows = servicesList.map((service: Service) => serviceRow(service));

  // Mock Data
  const data = {
    headers: ["Service", "Interaction", "Network", "Address", ""],
    rows: rows,
  };

  return <Table data={data} />;
};
