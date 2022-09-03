import {
  Address,
  DeleteDropdownItem,
  EditDropdownItem,
  InteractionBadge,
  NetworkBadge,
  OptionsDropdown,
} from "../../components/utils/Common";
import { getPage } from "../../constants/globals";
import Table from "../../components/table/Table";
import {
  services,
  Service,
  serviceFormFieldErrors,
  servicesFilter,
} from "../../recoil/atoms/services";
import { useRecoilState, useSetRecoilState } from "recoil";
import { removeFromList } from "../../recoil/atoms/common";
import { ServiceFormFields } from "./ServiceForm";
import { SimpleForm } from "../../components/forms/Form";

const EditService = ({ service }: { service: Service }) => {
  const pageName = "Services";

  // Get the page to set the icon
  const page = getPage(pageName);

  // Get the current values of the profile
  const setter = useSetRecoilState(servicesFilter(service.id));

  // Create a setter for the submit
  const onSubmit = (newElement: any) => {
    setter(newElement);
  };

  // Create the form with the default values as they currently are
  const content = (
    <SimpleForm
      create={false}
      defaultValues={service}
      errors={serviceFormFieldErrors}
      onSubmit={onSubmit}
      page={pageName}
      fields={ServiceFormFields}
    />
  );

  // Get the form with the update tag
  return (
    <EditDropdownItem form={content} icon={page?.icon} title={"profile"} />
  );
};

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
      <EditService service={service} />
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

export const ServicesTable = ({
  servicesList,
}: {
  servicesList: Service[];
}) => {
  const rows = servicesList.map((service: Service) => serviceRow(service));

  // Mock Data
  const data = {
    headers: ["Service", "Interaction", "Network", "Address", ""],
    rows: rows,
  };

  return <Table data={data} />;
};
