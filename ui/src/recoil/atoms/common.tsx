import uuid from "react-uuid";
import { Resetter, useRecoilState, useResetRecoilState } from "recoil";

export interface ValidatorProps {
  element: any;
  query: { [key: string]: any }[];
  errors: any;
}

export function ValidateIsFieldUniqueValue(
  { element, query, errors }: ValidatorProps,
  field: string
) {
  // Get the field value from the element
  const value = element[field];

  for (const q of query) {
    if (q[field] === value) {
      // Add the error to the list
      errors((prev: any) => ({ ...prev, [field]: "Invalid value" }));
    }
  }
}

export type ElementState = {
  getter: any;
  setter: any;
  resetter: Resetter;
};

export type AddWithValidatorsProps = {
  newElement: any;
  elementState: ElementState;
  errorState: ElementState;
  validators: any[];
};

export const AddWithValidators = ({
  newElement,
  elementState,
  errorState,
  validators,
}: AddWithValidatorsProps) => {
  // Reset the errors
  errorState.resetter();

  // Validate the fields
  const validatorProps = {
    element: newElement,
    query: elementState.getter,
    errors: errorState.setter,
  };
  validators.forEach((validate) => {
    validate(validatorProps);
  });

  // If errors has any key, return
  if (!!Object.keys(errorState.getter).length) return;

  // Add a new id
  let element: any = {
    ...newElement,
    id: uuid(),
  };

  // Set the new value in the list
  elementState.setter((oldList: any[]) => [...oldList, element]);

  // Reset everything
  [errorState, elementState].forEach((e) => e.resetter());
};

export const removeFromList = (id: string, getter: any[], setter: any) => {
  setter(getter.filter((service: any) => service.id !== id));
};

export type CreateItemProps = {
  filter: any;
  items: any;
  errors: any;
  validators: any[];
};

export const CreateItem = ({
  items,
  filter,
  errors,
  validators,
}: CreateItemProps) => {
  const [errorsGetter, errorsSetter] = useRecoilState(errors);
  const [getter, setter] = useRecoilState(items);

  // Resetters
  const resetFields = useResetRecoilState(filter);
  const resetErrors = useResetRecoilState(errors);

  return {
    elementState: {
      getter: getter,
      setter: setter,
      resetter: resetFields,
    },
    errorState: {
      getter: errorsGetter,
      setter: errorsSetter,
      resetter: resetErrors,
    },
    validators: validators,
  };
};
