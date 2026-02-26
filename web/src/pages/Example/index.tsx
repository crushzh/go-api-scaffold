import { PlusOutlined } from '@ant-design/icons';
import {
  ActionType,
  ModalForm,
  PageContainer,
  ProFormText,
  ProFormTextArea,
  ProFormSelect,
  ProTable,
  ProColumns,
} from '@ant-design/pro-components';
import { Button, message, Popconfirm } from 'antd';
import { useRef, useState } from 'react';
import {
  getExamples,
  createExample,
  updateExample,
  deleteExample,
} from '@/services/example';

type ExampleItem = {
  id: number;
  name: string;
  description: string;
  status: string;
  created_at: string;
  updated_at: string;
};

const ExamplePage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [modalOpen, setModalOpen] = useState(false);
  const [editRecord, setEditRecord] = useState<ExampleItem | null>(null);

  const columns: ProColumns<ExampleItem>[] = [
    { title: 'ID', dataIndex: 'id', width: 60, search: false },
    { title: 'Name', dataIndex: 'name', ellipsis: true },
    { title: 'Description', dataIndex: 'description', ellipsis: true, search: false },
    {
      title: 'Status',
      dataIndex: 'status',
      valueEnum: {
        active: { text: 'Active', status: 'Success' },
        inactive: { text: 'Inactive', status: 'Default' },
      },
    },
    {
      title: 'Created',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      search: false,
      width: 180,
    },
    {
      title: 'Action',
      valueType: 'option',
      width: 150,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setEditRecord(record);
            setModalOpen(true);
          }}
        >
          Edit
        </a>,
        <Popconfirm
          key="delete"
          title="Are you sure?"
          onConfirm={async () => {
            await deleteExample(record.id);
            message.success('Deleted');
            actionRef.current?.reload();
          }}
        >
          <a style={{ color: '#ff4d4f' }}>Delete</a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<ExampleItem>
        headerTitle="Example List"
        actionRef={actionRef}
        rowKey="id"
        columns={columns}
        search={{ labelWidth: 'auto' }}
        toolBarRender={() => [
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              setEditRecord(null);
              setModalOpen(true);
            }}
          >
            New
          </Button>,
        ]}
        request={async (params) => {
          const res = await getExamples({
            page: params.current,
            page_size: params.pageSize,
            keyword: params.name,
          });
          return {
            data: res?.data?.list || [],
            total: res?.data?.total || 0,
            success: true,
          };
        }}
        pagination={{ defaultPageSize: 10 }}
      />

      <ModalForm
        title={editRecord ? 'Edit Example' : 'New Example'}
        open={modalOpen}
        onOpenChange={setModalOpen}
        initialValues={editRecord || {}}
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editRecord) {
            await updateExample(editRecord.id, values);
            message.success('Updated');
          } else {
            await createExample(values);
            message.success('Created');
          }
          actionRef.current?.reload();
          return true;
        }}
      >
        <ProFormText
          name="name"
          label="Name"
          rules={[{ required: true, message: 'Please enter name' }]}
          placeholder="Enter name"
        />
        <ProFormTextArea
          name="description"
          label="Description"
          placeholder="Enter description"
        />
        {editRecord && (
          <ProFormSelect
            name="status"
            label="Status"
            options={[
              { label: 'Active', value: 'active' },
              { label: 'Inactive', value: 'inactive' },
            ]}
          />
        )}
      </ModalForm>
    </PageContainer>
  );
};

export default ExamplePage;
