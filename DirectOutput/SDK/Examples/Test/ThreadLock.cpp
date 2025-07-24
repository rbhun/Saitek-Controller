#include "stdafx.h"
#include "ThreadLock.h"

CThreadLock::CThreadLock()
{
	//InitializeCriticalSection(&m_sect);
	m_mutex = CreateMutex(NULL, FALSE, NULL);
}
CThreadLock::~CThreadLock()
{
	//DeleteCriticalSection(&m_sect);
	CloseHandle(m_mutex);
}
void CThreadLock::Acquire()
{
	//EnterCriticalSection(&m_sect);
	WaitForSingleObject(m_mutex, INFINITE);
}
void CThreadLock::Release()
{
	//LeaveCriticalSection(&m_sect);
	ReleaseMutex(m_mutex);
}
