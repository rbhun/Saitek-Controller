#pragma once

class CThreadLock
{
public:
	CThreadLock();
	~CThreadLock();

	void Acquire();
	void Release();
private:
	//CRITICAL_SECTION m_sect;
	HANDLE m_mutex;
};