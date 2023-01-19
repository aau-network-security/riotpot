import {
  Children,
  CSSProperties,
  forwardRef,
  ReactNode,
  Ref,
  useState,
} from "react";
import {
  Row,
  Col,
  Dropdown,
  Form,
  FormControl,
  ListGroup,
} from "react-bootstrap";
import { AiOutlinePlus } from "react-icons/ai";
import { useRecoilCallback, useRecoilState, useRecoilValue } from "recoil";
import { CustomToggle } from "../../components/dropdown/Dropdown";
import { InteractionBadge, NetworkBadge } from "../../components/utils/Common";
import { UtilsBar } from "../../components/utils/Utils";
import { profilesFilter } from "../../recoil/atoms/profiles";
import { Service, services } from "../../recoil/atoms/services";

type serviceHandler = {
  add: (service: Service) => void;
  remove: (service: Service) => void;
};

const ServiceCheckRow = ({
  service,
  checked,
  handler,
}: {
  service: Service;
  checked: boolean;
  handler: serviceHandler;
}) => {
  const [ch, setCh] = useState(checked);
  const { add, remove } = handler;

  const onChange = (e: any) => {
    setCh(e.target.checked);

    if (!e.target.checked) {
      remove(service);
      return;
    }

    add(service);
  };

  return (
    <Dropdown.Item key={service.id + "_dropdown"}>
      <Form.Check
        type="checkbox"
        onClick={(e) => e.stopPropagation()}
        onChange={onChange}
        checked={ch}
      />
      <Row>
        <Col className="service">{service.name}</Col>
        <Col className="badges">
          <NetworkBadge {...service.network} />
          <InteractionBadge {...service.interaction} />
        </Col>
      </Row>
    </Dropdown.Item>
  );
};

type CustomMenuProps = {
  children?: ReactNode;
  style?: CSSProperties;
  className?: string;
  labeledBy?: string;
};

const ProfileAddServiceDropdownMenu = forwardRef(
  (props: CustomMenuProps, ref: Ref<HTMLDivElement>) => {
    const [value, setValue] = useState("");

    return (
      <div
        ref={ref}
        style={props.style}
        className={props.className}
        aria-labelledby={props.labeledBy}
      >
        <FormControl
          autoFocus
          className="mx-3 my-2 w-auto"
          placeholder="Type to filter..."
          onChange={(e) => setValue(e.target.value)}
          value={value}
        />
        <ListGroup className="list-unstyled">
          {Children.toArray(props.children).filter(
            (child: any) =>
              !value || child.props.service.name.toLowerCase().startsWith(value)
          )}
        </ListGroup>
      </div>
    );
  }
);

const ServiceDropdownItems = ({ id }: { id: string }) => {
  const [profile, setProfile] = useRecoilState(profilesFilter(id));
  const servs = useRecoilValue(services);

  const insertService = useRecoilCallback(() => (service: Service) => {
    // Copy the original content
    var cp = { ...profile };
    let serviceIDs = cp.services.map((x: Service) => x.id);

    if (!serviceIDs.includes(service.id)) {
      cp.services = [...cp.services, service];
      // Place the profile with the updated values
      setProfile(cp);
    }
  });

  const removeService = useRecoilCallback(() => (service: Service) => {
    // Copy the original content
    let cp = { ...profile };
    cp.services = cp.services.filter((serv: Service) => serv.id !== service.id);

    // Place the profile with the updated values
    setProfile(cp);
  });

  const handler = {
    add: insertService,
    remove: removeService,
  };

  return (
    <Dropdown drop="start">
      <Dropdown.Toggle drop="start" as={CustomToggle} id={`dropdown-row-add`}>
        <AiOutlinePlus />
      </Dropdown.Toggle>
      <Dropdown.Menu as={ProfileAddServiceDropdownMenu}>
        {servs.map((serv: Service) => (
          <ServiceCheckRow
            key={serv.id}
            service={serv}
            handler={handler}
            checked={!!profile.services.find((s: Service) => s.id === serv.id)}
          />
        ))}
      </Dropdown.Menu>
    </Dropdown>
  );
};

export const ProfileUtils = ({ id }: { id: string }) => {
  const buttons = (
    <>
      <ServiceDropdownItems id={id} />
    </>
  );

  return <UtilsBar buttons={buttons} />;
};
