import React from "react";
import { Button } from "react-bootstrap";
import { AiOutlinePlus } from "react-icons/ai";
import { CenteredModal } from "../modal/Modal";
import "./Utils.scss";

export const Submit = () => {
  const onclick = () => {};

  return (
    <Button onClick={onclick} variant="success">
      Submit
    </Button>
  );
};

type CreateButtonProps = {
  icon?: any;
  title: string;
  content: any;
};

const CreateButton = ({ icon, title, content }: CreateButtonProps) => {
  const [modalShow, setModalShow] = React.useState(false);

  const props = {
    title: "New " + title,
    icon: icon,
    onHide: () => setModalShow(false),
    show: modalShow,
    submit: <Submit />,
    content: content,
  };

  return (
    <>
      <AiOutlinePlus onClick={() => setModalShow(true)}></AiOutlinePlus>
      <CenteredModal props={props} />
    </>
  );
};

/**
 *This component creates a utils bar with access to transform the current view.
 *
 */
const UtilsBar = ({ buttons = [] }: { buttons?: any[] }) => {
  return (
    <div className="utils">
      {buttons.map((btn) => {
        return btn;
      })}
    </div>
  );
};

export { UtilsBar, CreateButton };
