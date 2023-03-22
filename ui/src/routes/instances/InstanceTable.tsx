import { useEffect, useRef, useState } from "react";
import {
  Button,
  ButtonGroup,
  Col,
  Form,
  InputGroup,
  Row,
} from "react-bootstrap";
import { AiOutlineInfoCircle } from "react-icons/ai";
import { BsArrowRepeat, BsCheck, BsX } from "react-icons/bs";
import { FaNetworkWired } from "react-icons/fa";
import { useRecoilState, useRecoilValue } from "recoil";
import { Pop } from "../../components/pop/Pop";
import { Table, TableRow } from "../../components/table/Table";
import { useToast } from "../../components/toast/Toast";
import {
  DeleteDropdownItem,
  InteractionBadge,
  NetworkBadge,
  OptionsDropdown,
} from "../../components/utils/Common";
import { getPage } from "../../constants/globals";
import {
  InstanceProxy,
  GetInstanceAddress,
  useInstanceProxy,
  instanceProxySelector,
  proxies,
} from "../../recoil/atoms/instances";
import { Service } from "../../recoil/atoms/services";
import { ErrorToastVariant } from "../../recoil/atoms/toast";
import {
  changeProxyStatus,
  changeProxyPort,
  deleteProxyService,
} from "./InstanceAPI";
import InstanceUtils from "./InstanceUtils";

// Button to delete a proxy from the instance
const DeleteProxyServiceButton = ({ proxy }: { proxy: InstanceProxy }) => {
  const { removeProxy } = useInstanceProxy();

  const page = getPage("Services");
  const note = "The service will be stopped and removed from the instance";
  const address = GetInstanceAddress();

  const deleteCallback = () => {
    const deleted = deleteProxyService(address, proxy.id);
    deleted
      .then((data) => {
        if ("success" in data) {
          removeProxy(proxy);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  };

  return (
    <>
      {page && (
        <DeleteDropdownItem
          page={page}
          note={note}
          name={proxy.service.name}
          onClick={() => deleteCallback()}
        />
      )}
    </>
  );
};

// Dropdown with options to manage a proxy in an instance
const ProxyServiceRowOptions = ({ proxy }: { proxy: InstanceProxy }) => {
  return (
    <OptionsDropdown>
      <DeleteProxyServiceButton proxy={proxy} />
    </OptionsDropdown>
  );
};

// Shows service information badges about some service
const InstanceServiceInfo = ({ service }: { service: Service }) => {
  const [show, setShow] = useState(false);
  const target = useRef(null);

  return (
    <span>
      {service.name}
      <small className="info" ref={target}>
        <AiOutlineInfoCircle
          onMouseEnter={() => setShow(true)}
          onMouseLeave={() => setShow(false)}
        />
      </small>
      <Pop target={target} show={show} placement="left">
        <Row>
          <Col className="badges">
            <NetworkBadge {...service.network} />
            <InteractionBadge {...service.interaction} />
          </Col>
        </Row>
      </Pop>
    </span>
  );
};

// Section of the row that let you manipulate the proxy address
const InstanceServiceProxyAddress = ({ proxy }: { proxy: InstanceProxy }) => {
  const { showToast } = useToast();

  const [getProxy, setProxy] = useRecoilState(instanceProxySelector(proxy.id));
  const [getProxyPort, setProxyPort] = useState(getProxy?.port || 0);

  const address = GetInstanceAddress();

  const handler = () => {
    const proxyPort = changeProxyPort(address, proxy.id, getProxyPort);
    proxyPort
      .then((data) => {
        if ("error" in data) {
          showToast(data["error"], ErrorToastVariant);
          return;
        }
        setProxy(data);
      })
      .catch((error) => {
        showToast(error.message, ErrorToastVariant);
      });
  };

  return (
    <Col xs="8">
      <InputGroup size="sm" className="address">
        <InputGroup.Text>
          <FaNetworkWired />
        </InputGroup.Text>
        <Form.Control
          type="number"
          min={1}
          max={65535}
          defaultValue={getProxyPort}
          onChange={(e) => setProxyPort(parseInt(e.target.value))}
        />
        <Button
          variant="outline-secondary"
          id="button-addon2"
          onClick={() => handler()}
        >
          <BsArrowRepeat />
        </Button>
      </InputGroup>
    </Col>
  );
};

// Button to start or stop some proxy
const InstanceServiceToggle = ({ proxy }: { proxy: InstanceProxy }) => {
  const { showToast } = useToast();
  const [status, setStatus] = useState(proxy.status);
  const address = GetInstanceAddress();

  const handler = (isRunning: string) => {
    // Change the status of the thing
    if (status !== isRunning) {
      const response = changeProxyStatus(address, proxy.id, isRunning);

      response
        .then((data) => {
          if ("error" in data) {
            showToast(data["error"], ErrorToastVariant);
            return;
          }

          if (["running", "stopped"].includes(data.status))
            setStatus(data.status);
        })
        .catch((err) => {
          showToast(err.message, ErrorToastVariant);
        });
    }
  };

  return (
    <ButtonGroup aria-label="toggle-proxy" size="sm">
      <Button
        className="service-running running-true"
        variant="secondary"
        active={status === "running"}
        onClick={() => handler("running")}
      >
        <BsCheck />
      </Button>
      <Button
        className="service-running running-false"
        variant="secondary"
        active={status === "stopped"}
        onClick={() => handler("stopped")}
      >
        <BsX />
      </Button>
    </ButtonGroup>
  );
};

// Row that combines all the proxy service information
const InstanceServiceRow = ({ proxy }: { proxy: InstanceProxy }) => {
  const cells = [
    <InstanceServiceInfo key={0} service={proxy.service} />,
    <InstanceServiceProxyAddress key={1} proxy={proxy} />,
    <InstanceServiceToggle key={2} proxy={proxy} />,
    <ProxyServiceRowOptions key={3} proxy={proxy} />,
  ];

  return <TableRow cells={cells} />;
};

// Returns the actual table with all the rows
const InstanceServicesTable = () => {
  // Get all the proxy services set and create a row for each of them
  const pxs = useRecoilValue(proxies);

  // Map the rows into a proxy service
  const rows = pxs.map((proxy: InstanceProxy) => (
    <InstanceServiceRow key={proxy.id} proxy={proxy} />
  ));

  // Send the data
  const data = {
    headers: [`${pxs.length} Services`, "", "", ""],
    rows: [],
  };

  return (
    <>
      <InstanceUtils />
      <Table data={data} rows={rows}></Table>
    </>
  );
};

export default InstanceServicesTable;
