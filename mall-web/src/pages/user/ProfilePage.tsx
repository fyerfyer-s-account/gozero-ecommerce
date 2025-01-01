import React from 'react';
import { Card, Tabs } from 'antd';
import { useAuth } from '../../hooks/useAuth';
import ProfileForm from '../../components/user/ProfileForm';
import { Loading } from '../../components/common/Loading';

const { TabPane } = Tabs;

const ProfilePage: React.FC = () => {
  const { user, loading } = useAuth();

  if (loading) return <Loading />;

  if (!user) {
    return <div>Please login to view your profile.</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <Card className="max-w-2xl mx-auto">
        <Tabs defaultActiveKey="profile">
          <TabPane tab="Profile" key="profile">
            <h2 className="text-2xl font-bold mb-6">My Profile</h2>
            <ProfileForm />
          </TabPane>
          <TabPane tab="Security" key="security">
            <h2 className="text-2xl font-bold mb-6">Security Settings</h2>
            {/* Add security settings form here */}
          </TabPane>
          <TabPane tab="Preferences" key="preferences">
            <h2 className="text-2xl font-bold mb-6">Preferences</h2>
            {/* Add preferences settings here */}
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default ProfilePage;