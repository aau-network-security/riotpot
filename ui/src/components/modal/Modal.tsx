import "./Modal.scss";

import Modal from "react-bootstrap/Modal";

type ModalProps = {
  title: string;
  content: any;
  icon: any;

  // Modal state
  show: boolean;
  onHide: () => void;
};

/**
 * This component is used for regular modals
 * {title, icon, content} : {title: string, icon: any, content: any}
 */
const CenteredModal = ({
  props,
  children,
}: {
  props: ModalProps;
  children?: any;
}) => {
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
      <Modal.Body>
        {children}
        {props.content}
      </Modal.Body>
      <Modal.Footer></Modal.Footer>
    </Modal>
  );
};

export { CenteredModal };
