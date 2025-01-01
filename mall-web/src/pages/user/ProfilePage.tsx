import React from 'react';
import { Card, Descriptions, Tabs } from 'antd';
import { useAuth } from '@/hooks/useAuth';
import ProfileForm from '@/components/user/ProfileForm';
import { Loading } from '@/components/common/Loading';

const ProfilePage: React.FC = () => {
  const { user, loading } = useAuth();

  if (loading) return <Loading />;

  if (!user) {
    return <div>Please login to view your profile.</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <Card className="max-w-2xl mx-auto">
        <Tabs defaultActiveKey="info">
          <Tabs.TabPane tab="Basic Info" key="info">
            <Descriptions bordered column={1}>
              <Descriptions.Item label="Username">{user.username}</Descriptions.Item>
              <Descriptions.Item label="Member Level">Level {user.memberLevel}</Descriptions.Item>
              <Descriptions.Item label="Balance">${user.balance}</Descriptions.Item>
              <Descriptions.Item label="Join Date">
                {new Date(user.createdAt * 1000).toLocaleDateString()}
              </Descriptions.Item>
            </Descriptions>
          </Tabs.TabPane>
          
          <Tabs.TabPane tab="Edit Profile" key="edit">
            <ProfileForm />
          </Tabs.TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default ProfilePage;