import { useRef, useState } from "react";
import { Button, Col, Dropdown, Form, InputGroup } from "react-bootstrap";
import {
  AiFillCheckCircle,
  AiFillCloseCircle,
  AiOutlineInfoCircle,
} from "react-icons/ai";
import { BsArrowRepeat, BsGlobe } from "react-icons/bs";
import { CgDetailsLess } from "react-icons/cg";
import { FaNetworkWired } from "react-icons/fa";
import { Link } from "react-router-dom";
import { useRecoilCallback, useRecoilState, useRecoilValue } from "recoil";
import { SimpleForm } from "../../components/forms/Form";
import { Pop } from "../../components/pop/Pop";
import { Table, Row } from "../../components/table/Table";
import {
  DeleteDropdownItem,
  EditDropdownItem,
  InteractionBadge,
  OptionsDropdown,
} from "../../components/utils/Common";
import { getPage, InteractionOption } from "../../constants/globals";
import {
  Instance,
  instanceIds,
  instances,
  intanceFormFieldErrors,
} from "../../recoil/atoms/instances";
import { InstanceFormFields } from "./InstanceForm";

interface InstanceService {
  name: string;
  proxy: Number;
  interaction: InteractionOption;
  running: boolean;
}

const EditInstance = ({ instance }: { instance: Instance }) => {
  const pageName = "Instances";
  // Get the page to set the icon
  const page = getPage(pageName);

  // Create a setter for the submit
  const id = instance.id !== undefined ? instance.id : -1;
  const onSubmit = useRecoilCallback(({ set }) => (instance: Instance) => {
    set(instances(id), (prev) => ({ ...prev, ...instance }));
  });

  // Create the form with the default values as they currently are
  const content = (
    <SimpleForm
      create={false}
      defaultValues={instance}
      errors={intanceFormFieldErrors}
      onSubmit={onSubmit}
      page={pageName}
      fields={InstanceFormFields}
    />
  );

  // Get the form with the update tag
  return (
    <EditDropdownItem form={content} icon={page?.icon} title={"profile"} />
  );
};

const ViewInstance = ({ instance }: { instance: Instance }) => {
  return (
    <Dropdown.Item>
      <Link to={`${instance.id}`}>
        <CgDetailsLess />
        Details
      </Link>
    </Dropdown.Item>
  );
};

const ProxyInfoPop = ({ proxy }: { proxy: Number }) => {
  return (
    <span className="proxy-info">
      <FaNetworkWired />
      {`${proxy}`}
    </span>
  );
};

const InstancePop = ({ service }: { service: InstanceService }) => {
  const cross = <AiFillCloseCircle style={{ fill: "#FF8686" }} />;
  const tick = <AiFillCheckCircle style={{ fill: "#A4E18F" }} />;

  return (
    <Col>
      {service.running ? tick : cross} {service.name}
      <ProxyInfoPop proxy={service.proxy} />
      <InteractionBadge {...service.interaction} />
    </Col>
  );
};

const InstanceRowInfoPop = ({ services }: { services: InstanceService[] }) => {
  const [show, setShow] = useState(false);
  const target = useRef(null);

  return (
    <>
      <small className="info" ref={target} onClick={() => setShow(!show)}>
        <AiOutlineInfoCircle />
        {`${services.length} ${services.length === 1 ? "service" : "services"}`}
      </small>
      <Pop target={target} show={show} placement="right">
        {services.map((service) => {
          return <InstancePop service={service} />;
        })}
      </Pop>
    </>
  );
};

const InstanceRowInfo = ({
  name,
  services,
}: {
  name: string;
  services: InstanceService[];
}) => {
  return (
    <>
      <div>{name}</div>
      <InstanceRowInfoPop services={services} />
    </>
  );
};

const InstanceRowAddress = ({ id }: { id: number }) => {
  const [instance, setter] = useRecoilState(instances(id));
  const [value, setValue] = useState(instance.host);
  const onClick = () => {
    setter({ ...instance, host: value });
  };

  return (
    <Col xs="8">
      <InputGroup size="sm" className="address">
        <InputGroup.Text>
          <BsGlobe />
        </InputGroup.Text>
        <Form.Control
          placeholder="Address"
          aria-label="Instance's address"
          aria-describedby="basic-addon2"
          value={value}
          onChange={(e) => setValue(e.target.value)}
        />
        <Button
          variant="outline-secondary"
          id="button-addon2"
          onClick={onClick}
        >
          <BsArrowRepeat />
        </Button>
      </InputGroup>
    </Col>
  );
};

const InstanceRowOptions = ({ instance }: { instance: Instance }) => {
  const ids = useRecoilValue(instanceIds);

  const deleteCallback = useRecoilCallback(
    ({ set }) =>
      (id: number | undefined) => {
        set(
          instanceIds,
          ids.filter((curr) => curr !== id)
        );
      }
  );

  const page = getPage("Instances");
  const note =
    "Instance services will be stopped and removed from the instance register.";
  return (
    <OptionsDropdown>
      <ViewInstance instance={instance} />
      <EditInstance instance={instance} />
      {page && (
        <DeleteDropdownItem
          page={page}
          note={note}
          name={instance.name}
          onClick={() => deleteCallback(instance.id)}
        />
      )}
    </OptionsDropdown>
  );
};

const InstanceRow = ({ id }: { id: number }) => {
  const instance = useRecoilValue(instances(id));
  const services: InstanceService[] = [];

  const cells = [
    <InstanceRowInfo name={instance.name} services={services} />,
    <InstanceRowAddress id={id} />,
    <InstanceRowOptions instance={instance} />,
  ];

  return <Row cells={cells} />;
};

export const InstancesTable = () => {
  const insts = useRecoilValue(instanceIds);
  const rows = insts.map((instance, index) => (
    <InstanceRow key={index} id={instance} />
  ));

  const data = {
    headers: [`${rows.length} Instances`, "", ""],
    rows: [],
  };

  return <Table data={data} rows={rows} />;
};
