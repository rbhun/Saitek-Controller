#pragma once
#include "ThreadLock.h"

class CRawImage
{
public:
	CRawImage(int width = 320, int height = 240, int bpp = 24);
	~CRawImage();

	HDC  BeginPaint();
	void EndPaint();

	void Acquire();
	void Release();

	// protected by Acquire() / Release()
	DWORD Size() const;
	// protected by Acquire() / Release()
	const void* Buffer() const;
private:
	CThreadLock m_lock;

	HDC			m_hdc;
	HBITMAP		m_bmp;
	HBITMAP		m_old;
	BITMAPINFO	m_info;

	LPBYTE		m_raw;
	LPBYTE		m_return;
};