import { Toast, ToastContainer } from "react-bootstrap";
import { useRecoilState } from "recoil";
import {
  InfoToast,
  InfoToastVariant,
  toast,
  ToastState,
  ToastVariant,
} from "../../recoil/atoms/toast";

export function useToast() {
  const [t, setToast] = useRecoilState<ToastState>(toast);

  const showToast = (message: string, variant?: ToastVariant) => {
    // Use the default version of an info toast
    let def = InfoToast;

    let upd = {
      ...def,
      variant: variant || InfoToastVariant,
      message: message,
      show: true,
    };

    setToast(upd);
  };

  return { t, showToast };
}

export const CustomToast = () => {
  const [t, setToast] = useRecoilState<ToastState>(toast);

  const hideToast = () => {
    setToast(InfoToast);
  };

  return (
    <ToastContainer className="p-3" position="bottom-end">
      <Toast
        bg={t.variant.color}
        onClose={hideToast}
        show={t.show}
        delay={3000}
        autohide
      >
        <Toast.Header>
          {t.variant.icon}
          <strong className="me-auto">{t.variant.name}</strong>
        </Toast.Header>
        <Toast.Body>{t.message}</Toast.Body>
      </Toast>
    </ToastContainer>
  );
};
