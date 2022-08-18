import { SimpleForm } from "../../components/forms/Form";

type InstanceFormProps = {
  name: string;
  description: string;
  host: string;
};

const InstanceForm = ({ data }: { data?: InstanceFormProps }) => {
  // TODO: Perhaps repalce this in the future with an API call to get the fields and their structure
  let fields = [
    {
      name: "name",
      help: "Name of the instance",
      fieldattrs: {
        type: "text",
        required: true,
        value: data?.name,
        placeholder: "Type a name for the instance",
      },
    },
    {
      name: "description",
      help: "Instance description",
      fieldattrs: {
        type: "text",
        required: false,
        value: data?.description,
        placeholder: "Add a short description for the instance",
        as: "textarea",
      },
    },
    {
      name: "host",
      help: "Host in where the instance can be reached. E.g.: localhost, 0.0.0.0, etc.",
      fieldattrs: {
        type: "text",
        required: true,
        value: data?.host,
        placeholder: "localhost",
      },
    },
  ];

  return <SimpleForm page="Instances" fields={fields}></SimpleForm>;
};

export default InstanceForm;
