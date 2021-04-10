import React, { useEffect, useState } from 'react';
import { ComponentGeneratorFunc } from '../../../../../types/generators';
import { Field, Resource, Schema } from '../../../../../types/schema';
// import { guessMainTitleField } from '../../../../utils/field';
import { DataGrid } from '@material-ui/data-grid';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { useDataProvider, useRedirect } from 'ra-core';
import { Link } from 'react-router-dom';


interface DatagridColumn {
  field: string;
  headerName: string;
  flex: number;
  width?: number;
}

interface ReferenceListDatagridProps {
  resource: Resource;
  refResourceFieldName: string;
  recordID: string;
}

const ReferenceListDatagrid = (props: ReferenceListDatagridProps) => {
  const recordID = props.recordID;
  const redirect = useRedirect();
  const [rows, setRows] = useState<any[]>([]);
  const dp = useDataProvider();
  useEffect(() => {
    const filter: any = {};
    filter[`${props.refResourceFieldName}_eq`] = recordID;
    dp.getList(props.resource.endpoint, {
      filter,
      pagination: {
        page: 1,
        perPage: 10,
      },
      sort: {
        field: 'id',
        order: "ASC",
      }
    }).then(res => {
      setRows(res?.data || []);
    })
  }, [dp, props.refResourceFieldName, props.resource.endpoint, recordID, setRows]);
  const columns: DatagridColumn[] = props.resource.list_endpoint.fields.map(f => ({
    field: f.name,
    headerName: f.options["label"],
    flex: f.name === "id" ? 0.3 : 1,
  })); 
  return(
    <div style={{ height: 400, width: '100%', marginBottom: '32px' }}>
      <DataGrid
        rows={rows}
        columns={columns}
        pageSize={5}
        rowHeight={48}
        disableColumnMenu
        disableSelectionOnClick
        disableColumnSelector
        componentsProps={{
          toolbar: {
            resource: props.resource,
            refResourceFieldName: props.refResourceFieldName,
            recordID: recordID,
          },
        }}
        onRowClick={(r) => {
          redirect( `/${props.resource.endpoint}/${r.row.id}`);
        }}
        components={{
          Toolbar: (props: any) => {
            const record: any = {};
            record[props.refResourceFieldName] = props.recordID;
            return (
              <div style={{ padding: '20px', borderBottom: '1px solid #e9e9e9', overflow: 'auto' }}> 
                <div style={{ float: 'left' }}>
                  <Typography variant="h3" style={{ fontSize: '20px' }}>
                    {props.resource.name}
                  </Typography>
                </div>
                <div style={{ float: 'right' }}>             
                  <Button
                    size="small"
                    variant="contained"
                    color="primary"
                    component={Link}
                    to={{
                        pathname: `/${props.resource.endpoint}/create`,
                        state: { record: record },
                    }}
                  >
                    Create
                  </Button>
                </div>       
              </div>
            );
          }
        }}
      />
    </div>
  );
};

const ReferenceListDatagridGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => {
  const resource = s.admin.resources.find(r => r.name === f.options["ref_resource"]);
  if (!resource) {
    console.error(`[ReferenceInput] Resource not found "${resource}"`);
    return null;
  }
  const { ref_resource_field_name } = f.options;
  return (
    <ReferenceListDatagrid resource={resource} refResourceFieldName={ref_resource_field_name} recordID={props.recordID} />
  );
}


export const ReferenceListFieldGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => {
  const resource = s.admin.resources.find(r => r.name === f.options["ref_resource"]);
  if (!resource) {
    console.error(`[ReferenceInput] Resource not found "${resource}"`);
    return null;
  }
  // const optTextFieldName = f.options["text_field_name"];
  // const textFieldName = optTextFieldName && optTextFieldName !== "" ? optTextFieldName : guessMainTitleField(resource);
  return (
    <div>
      TODO: ReferenceListField
    </div>
  );
}

export const ReferenceListInputGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => {
  const Comp = ReferenceListDatagridGenerator(s, f);
  return (props: any) => {
    return (
      <Comp {...props} recordID={props.record.id} />
    );
  };
};
