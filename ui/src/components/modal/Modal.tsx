import "./Modal.scss";

import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";

type ModalProps = {
  title: string;
  icon?: any;
  content: any;
  show: boolean;
  onHide: () => void;
  submit?: any;
};

/**
 * This component is used for regular modals
 * {title, icon, content} : {title: string, icon: any, content: any}
 */
const CenteredModal = ({ props }: { props: ModalProps }) => {
  return (
    <Modal
      {...props}
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      centered
    >
      <Modal.Header closeButton closeVariant="white">
        <Modal.Title id="contained-modal-title-vcenter">
          <props.icon />
          {props.title}
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>{props.content}</Modal.Body>
      <Modal.Footer>
        <Button onClick={props.onHide} variant="outline-light">
          Close
        </Button>
        {props.submit}
      </Modal.Footer>
    </Modal>
  );
};

export { CenteredModal };
