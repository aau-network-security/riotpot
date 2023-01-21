import "./Table.scss";

export const Cell = ({ content }: { content: any }) => {
  return <td className="tableCell">{content}</td>;
};

const Header = ({ content }: { content: any }) => {
  return <th className="tableCell">{content}</th>;
};

const Headers = ({ headers }: { headers: any }) => {
  return (
    <thead>
      <tr className="tableHeaders">
        {headers.map((cell: any, ind: number) => {
          const key = `header-${ind}`;
          return <Header content={cell} key={key} />;
        })}
      </tr>
    </thead>
  );
};

export const TableRow = ({ cells }: { cells: any }) => {
  return (
    <tr className="tableRow">
      {cells.map((content: any, ind: number) => {
        const key = `cell-${ind}`;
        return <Cell content={content} key={key} />;
      })}
    </tr>
  );
};

const Body = ({ rows, children }: { rows: any[]; children: any[] }) => {
  return (
    <tbody>
      {rows.map((cells: any[], ind: number) => {
        const key = `tablerow-${ind}`;
        return <TableRow cells={cells} key={key} />;
      })}
      {children}
    </tbody>
  );
};

export const Table = ({ data, rows }: { data: any; rows?: any }) => {
  return (
    <div className="table-container">
      <table>
        <Headers headers={data.headers} />
        <Body rows={data.rows}>{rows}</Body>
      </table>
    </div>
  );
};

export default Table;
