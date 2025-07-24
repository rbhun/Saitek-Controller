#include "stdafx.h"
#include "RawImage.h"

CRawImage::CRawImage(int width /*= 320*/, int height /*= 240*/, int bpp /*= 24*/) : 
		m_old(NULL), m_raw(NULL), m_return(NULL)
{
	// initialize BMP data
	m_hdc = CreateCompatibleDC(NULL);

	memset(&m_info, 0, sizeof(m_info));

	m_info.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
	m_info.bmiHeader.biWidth = width;
	m_info.bmiHeader.biHeight = height;
	m_info.bmiHeader.biPlanes = 1;
	m_info.bmiHeader.biBitCount = (WORD)bpp;
	m_info.bmiHeader.biCompression = BI_RGB;
	m_info.bmiHeader.biSizeImage = (m_info.bmiHeader.biWidth * m_info.bmiHeader.biHeight * m_info.bmiHeader.biBitCount) / 8;
	m_info.bmiHeader.biXPelsPerMeter = 3200;
	m_info.bmiHeader.biYPelsPerMeter = 3200;
	m_info.bmiHeader.biClrImportant = 0;
	m_info.bmiHeader.biClrUsed = 0;

	m_bmp = CreateDIBSection(m_hdc, &m_info, DIB_RGB_COLORS, (PVOID*)&m_raw, NULL, 0);
	memset(m_raw, 0, m_info.bmiHeader.biSizeImage);

	m_return = new BYTE[m_info.bmiHeader.biSizeImage];
	if (m_return)
	{
		memset(m_return, 0, m_info.bmiHeader.biSizeImage);
	}
}
CRawImage::~CRawImage()
{
	DeleteObject(m_bmp);
	DeleteDC(m_hdc);
	if (m_return)	delete [] m_return;
}
HDC  CRawImage::BeginPaint()
{
	m_old = (HBITMAP)SelectObject(m_hdc, m_bmp);
	return m_hdc;
}
void CRawImage::EndPaint()
{
	// flush GDI
	GdiFlush();
	// unselect BMP from HDC
	SelectObject(m_hdc, m_old);
	m_old = NULL;
	// make a copy of the raw buffer into the return buffer
	m_lock.Acquire();
	if (m_return && m_raw)
	{
		memcpy(m_return, m_raw, m_info.bmiHeader.biSizeImage);
	}
	m_lock.Release();
}
void CRawImage::Acquire()
{
	m_lock.Acquire();
}
void CRawImage::Release()
{
	m_lock.Release();
}
// protected by Acquire() / Release()
DWORD CRawImage::Size() const
{
	return m_info.bmiHeader.biSizeImage;
}
// protected by Acquire() / Release()
const void* CRawImage::Buffer() const
{
	return m_return;
}
