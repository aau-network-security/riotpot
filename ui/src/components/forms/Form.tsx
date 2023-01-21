import { FormEvent } from "react";
import { Button, Form } from "react-bootstrap";
import Select from "react-select";
import { useRecoilValue } from "recoil";
import { customTheme } from "../../constants/globals";

export type FormFieldProps = {
  name: string;
  type?: string;
  fieldattrs: { [key: string]: any };
  value?: string | number | string[] | undefined;
  help?: string;
  readOnly?: boolean;
};

export const FormField = ({
  page,
  fieldProps,
  type,
  defaultValue,
  error,
}: {
  page: string;
  type?: string;
  fieldProps: FormFieldProps;
  defaultValue?: any;
  error?: any;
}) => {
  // Get the state value. We only use the value as default. The value will be set
  // once the form is submitted!
  const fieldError: any = useRecoilValue(error(fieldProps.name));

  let name = fieldProps.name;
  name = (name.charAt(0).toUpperCase() + name.slice(1)).replace("_", " ");

  let field;

  switch (type) {
    case "select":
      field = (
        <Select
          {...fieldProps.fieldattrs}
          name={fieldProps.name}
          className="basic-single"
          classNamePrefix="select"
          defaultValue={defaultValue}
          isSearchable={true}
          styles={customTheme}
        />
      );
      break;
    default:
      field = (
        <Form.Control
          {...fieldProps.fieldattrs}
          name={fieldProps.name}
          defaultValue={defaultValue}
          isInvalid={!!fieldError}
        />
      );
  }

  return (
    <Form.Group className="mb-3" controlId={`form${page}${name}`}>
      <Form.Label>{name}</Form.Label>
      {field}
      {fieldError && (
        <Form.Control.Feedback type="invalid" role="alert">
          Invalid value
        </Form.Control.Feedback>
      )}
      <Form.Text className="text-muted">{fieldProps.help}</Form.Text>
    </Form.Group>
  );
};

export interface FormProps {
  page: string;
  fields: FormFieldProps[];
  onSubmit: any;
  defaultValues?: { [key: string]: any };
  errors?: any;
  create?: boolean;
}

export const SimpleForm = ({
  page,
  fields,
  defaultValues,
  errors,
  onSubmit,
  create,
}: FormProps) => {
  const onsubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault(); // Stop the page from refreshing

    // Map the values of the fields array to an object
    const ret: any = {};
    fields.forEach((field) => {
      let value;

      const target = event.currentTarget[field.name];

      switch (field.type) {
        // If the field is a "select", find the option and set it as the value
        case "select":
          value = field.fieldattrs.options.find(
            (x: any) => x.value === target.value
          );
          break;
        default:
          value = target.value;
      }

      ret[field.name] = value;
    });

    // Add the new value to the onSubmit function
    onSubmit(ret);
  };

  // If the intention is to create a new object, remove readOnly fields
  if (create === false) {
    fields = fields.filter((field) => !field.readOnly);
  }

  return (
    <Form onSubmit={onsubmit}>
      {fields.map((field, ind) => {
        const key = `form-${ind}`;

        return (
          <FormField
            type={field.type}
            page={page}
            fieldProps={field}
            defaultValue={defaultValues && defaultValues[field.name]}
            error={errors}
            key={key}
          />
        );
      })}
      <Button variant="primary" type="submit">
        Submit
      </Button>
    </Form>
  );
};
