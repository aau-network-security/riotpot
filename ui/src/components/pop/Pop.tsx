import { Container, Overlay } from "react-bootstrap";
import { Placement } from "react-bootstrap/esm/types";

import "./Pop.scss";

type PopProps = {
  target: any;
  children: any;
  show: boolean;
  placement: Placement;
};

export const Pop = ({
  target,
  children,
  show,
  placement = "right",
}: PopProps) => {
  return (
    <Overlay target={target.current} show={show} placement={placement}>
      <div className="component">
        <Container>{children}</Container>
      </div>
    </Overlay>
  );
};
