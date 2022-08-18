import { Form } from "react-bootstrap";

export type FormFieldProps = {
  name: string;
  help?: string;
  fieldattrs: object;
};

export const FormField = ({
  page,
  field,
}: {
  page: string;
  field: FormFieldProps;
}) => {
  var name = field.name;
  name = (name.charAt(0).toUpperCase() + name.slice(1)).replace("_", " ");

  return (
    <Form.Group className="mb-3" controlId={`form${page}${name}`}>
      <Form.Label>{name}</Form.Label>
      <Form.Control {...field.fieldattrs}></Form.Control>
      <Form.Text className="text-muted">{field.help}</Form.Text>
    </Form.Group>
  );
};

export const SimpleForm = ({
  page,
  fields,
  children,
}: {
  page: string;
  fields?: FormFieldProps[];
  children?: any;
}) => {
  return (
    <Form>
      {fields &&
        fields.map((field) => {
          return <FormField page={page} field={field} />;
        })}
      {children}
    </Form>
  );
};
