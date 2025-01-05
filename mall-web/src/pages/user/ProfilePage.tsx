import React, { useEffect, useState } from 'react';
import { Card, Descriptions, Tabs, message } from 'antd';
import { useAuth } from '@/hooks/useAuth';
import ProfileForm from '@/components/user/ProfileForm';
import { Loading } from '@/components/common/Loading';
import { useNavigate } from 'react-router-dom';

const ProfilePage: React.FC = () => {
  const { user, loading } = useAuth();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('info');

  useEffect(() => {
    if (!loading && !user) {
      message.error('Please login first');
      navigate('/login');
    }
  }, [user, loading]);

  const userFields = [
    { label: "Username", value: user?.username },
    { label: "Nickname", value: user?.nickname },
    { label: "Email", value: user?.email },
    { label: "Phone", value: user?.phone },
    { label: "Gender", value: user?.gender },
    { label: "Member Level", value: user?.memberLevel ? `Level ${user.memberLevel}` : null },
    { label: "Balance", value: user?.balance ? `$${user.balance}` : null },
    { label: "Join Date", value: user?.createdAt ? new Date(user.createdAt * 1000).toLocaleDateString() : null }
  ].filter(field => field.value); // Filter out empty values

  return (
    <div className="container mx-auto px-4 py-8">
      <Card className="max-w-2xl mx-auto">
        <Tabs activeKey={activeTab} onChange={setActiveTab}>
          <Tabs.TabPane tab="Basic Info" key="info">
            <Descriptions bordered column={1}>
              {userFields.map(field => (
                <Descriptions.Item key={field.label} label={field.label}>
                  {field.value}
                </Descriptions.Item>
              ))}
            </Descriptions>
          </Tabs.TabPane>
          
          <Tabs.TabPane tab="Edit Profile" key="edit">
            <ProfileForm onSuccess={() => setActiveTab('info')} />
          </Tabs.TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default ProfilePage;