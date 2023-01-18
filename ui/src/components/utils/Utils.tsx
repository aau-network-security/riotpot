import { useState } from "react";
import { AiOutlinePlus } from "react-icons/ai";
import { CenteredModal } from "../modal/Modal";
import "./Utils.scss";

type CreateButtonProps = {
  icon?: any;
  title: string;
  content: any;
};

const CreateButton = ({ icon, title, content }: CreateButtonProps) => {
  const [modalShow, setModalShow] = useState(false);

  const props = {
    title: "New " + title,
    show: modalShow,
    content: content,
    icon: icon,
    onHide: () => setModalShow(false),
  };

  return (
    <>
      <AiOutlinePlus
        onClick={() => setModalShow(true)}
        style={{
          cursor: "pointer",
        }}
      ></AiOutlinePlus>
      <CenteredModal props={props} />
    </>
  );
};

/**
 *This component creates a utils bar with access to transform the current view.
 */
const UtilsBar = ({ buttons }: { buttons?: any }) => {
  return <div className="utils">{buttons}</div>;
};

export { UtilsBar, CreateButton };
