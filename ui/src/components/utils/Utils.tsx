import React from "react";
import { AiOutlinePlus } from "react-icons/ai";
import { CenteredModal } from "../modal/Modal";
import "./Utils.scss";

type CreateButtonProps = {
  icon?: any;
  title: string;
  content: any;
};

const CreateButton = ({ icon, title, content }: CreateButtonProps) => {
  const [modalShow, setModalShow] = React.useState(false);

  const props = {
    title: "New " + title,
    show: modalShow,
    content: content,
    icon: icon,
    onHide: () => setModalShow(false),
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
