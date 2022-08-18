import "./Table.scss";

const Cell = ({ content }: { content: any }) => {
  return <td className="tableCell">{content}</td>;
};

const Header = ({ content }: { content: any }) => {
  return <th className="tableCell">{content}</th>;
};

const Headers = ({ headers }: { headers: any }) => {
  return (
    <thead>
      <tr className="tableHeaders">
        {headers.map((cell: any) => {
          return <Header content={cell} />;
        })}
      </tr>
    </thead>
  );
};

const Row = ({ cells }: { cells: any }) => {
  return (
    <tr className="tableRow">
      {cells.map((content: any, ind: Number) => {
        return <Cell content={content} />;
      })}
    </tr>
  );
};

const Body = ({ rows }: { rows: any }) => {
  return (
    <tbody>
      {rows.map((cells: any, ind: Number) => {
        return <Row cells={cells} />;
      })}
    </tbody>
  );
};

const Table = ({ data }: { data: any }) => {
  return (
    <div className="table-container">
      <table>
        <Headers headers={data.headers} />
        <Body rows={data.rows} />
      </table>
    </div>
  );
};

export default Table;
