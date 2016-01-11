using System;
using System.Collections.Generic;
using System.Text;
using System.IO;
using System.Drawing;

namespace Heroes5RandomMap.DataProc
{
    public enum TerrainDataError
    {
        None = 0,
        TagMismatch, //标记不匹配
        NotEnoughByte, //长度不足
        TagZero, //该段开头为空
    }

    public interface IColorize<T> //把一个数字转换为颜色以显示在屏幕上
    {
        Color PointToColor(T obPoint);
    }
    public class IGlobalToColor<T> : IColorize<T>
    {
        public Color PointToColor(T obPoint)
        {
            return System.Drawing.Color.FromArgb(0, 0, 0);
        }
    }
    public class INumberToColor<T> : IColorize<T> where T : struct
    {
        public Color PointToColor(T obPoint)
        {
            int RGBValue = 0;
            return System.Drawing.Color.FromArgb(RGBValue, RGBValue, RGBValue);
        }
    }
    public class IByteToColor : IColorize<Byte>
    {
        public Color PointToColor(Byte Point)
        {
            return System.Drawing.Color.FromArgb((int)Point, (int)Point, (int)Point);
        }
    }
    public class IHeightToColor : IColorize<Single> //1-1024的浮点数，实际上一般不到50
    {
        public Color PointToColor(Single Point)
        {
            int RGBValue = (int)Math.Abs(Math.Floor(Point)) * 5;
            return System.Drawing.Color.FromArgb(RGBValue, RGBValue, RGBValue);
        }
    }
    public class IUInt64ToColor : IColorize<UInt64>
    {
        public Color PointToColor(UInt64 Point)
        {
            byte[] NumberByte = BitConverter.GetBytes(Point);
            int RGBValue = (int)NumberByte[7];
            return System.Drawing.Color.FromArgb(RGBValue, RGBValue, RGBValue);
        }
    }
    public class IPlateauToColor : IColorize<Byte> // 只有0,10，20，30数据
    {
        public Color PointToColor(Byte Point)
        {
            int RGBValue = (int)(Point) * 5;
            return System.Drawing.Color.FromArgb(RGBValue, RGBValue, RGBValue);
        }
    }
    public class IPassibleToColor : IColorize<Byte> // 只有0，1二种数据
    {
        public Color PointToColor(Byte Point)
        {
            int RGBValue = (int)Point * 0xFF;
            return System.Drawing.Color.FromArgb(RGBValue, RGBValue, RGBValue);
        }
    }

    public interface IHexViewable<T> // 这个接口用于把数据转换成十六进制格式显示在文本框里面，浮点例外，为一位小数
    {
        string ToHexString(T obNumber);
    }
    public abstract class ABaseToHexView<T> : IHexViewable<T>
    {
        public abstract string ToHexString(T obNumber);
    } //抽象类作为基类
    public class IGlobalToHexView<T> : ABaseToHexView<T>//全局通用的借口，返回"N/A"
    {
        override public string ToHexString(T obNumber)
        {
            return "N/A";
        }
    }
    public class IByteArrayToHexView : ABaseToHexView<byte[]>
    {
        override public string ToHexString(byte[] bArray)
        {
            return bArray[0].ToString("X2");
        }
    }
    public class INumberToHexView<T> : ABaseToHexView<T> where T : struct, IFormattable //用于数字通用的接口类，返回直接转为16进制的数
    {
        override public string ToHexString(T obNumber)
        {
            return obNumber.ToString("X2", null);
        }
    }
    public class IHeightToHexView : ABaseToHexView<Single>
    {
        override public string ToHexString(Single obNumber)
        {
            return ((int)Math.Abs(obNumber)).ToString("F1", null);
        }
    }
    public class IUInt64ToHexView : ABaseToHexView<UInt64>
    {
        override public string ToHexString(UInt64 obNumber)
        {
            byte[] NumberByte = BitConverter.GetBytes(obNumber);
            return NumberByte[7].ToString("X2", null);
        }
    }

    public interface IExportGridInfo // 这个接口用于输出Color和HexView，以及读入读出文件。TerrainDataBlock具有这个接口
    {
        string[] ExportText();
        Bitmap ExportBMP(int Width, int Height);
        void Save(BinaryWriter brData);
        TerrainDataError Load(BinaryReader brData, byte[] anBlockTag, bool ExistSizeTag);
    }

    public class TerrainDataBlock<T> : IExportGridInfo where T : struct, IFormattable
    {
        public T[,] tContent;
        public byte[] nXLengthTag = new byte[2] { 0x01, 0x08 };
        public byte[] nYLengthTag = new byte[2] { 0x02, 0x08 };
        public byte[] nBlockTag;
        public byte[] nSizeTag = new byte[1] { 0x03 }; // 注意这个变量取0的时候就是没有尺寸标记
        public UInt32 nSize; // 总长度
        public bool bExistTag = true;

        public IColorize<T> iColorize = new IGlobalToColor<T>();
        public IHexViewable<T> iHexViewer = new INumberToHexView<T>();
        //构造函数
        public TerrainDataBlock(byte[] BlockTag) //用数据标记初始化空类
        {
            nBlockTag = BlockTag;
        }
        public TerrainDataBlock(byte[] BlockTag, UInt32 Size):this(BlockTag) //用数据标记和数据块边长初始化空类，默认为方形
        {
            CreateContent(Size);
        }
        public TerrainDataBlock(byte[] BlockTag, IColorize<T> iPointToColor, IHexViewable<T> iPointToHex)
            : this(BlockTag) //用标志和两个外部接口初始化类
        {
            iColorize = iPointToColor;
            iHexViewer = iPointToHex;
        }
        public TerrainDataBlock(byte[] BlockTag, UInt32 Size, IColorize<T> iPointToColor, IHexViewable<T> iPointToHex):this(BlockTag,Size)
            //初始化所有的4个参数
        {
            iColorize = iPointToColor;
            iHexViewer = iPointToHex;
        }

        public void CreateContent(UInt32 nSize)//初始化tContent，用默认0值
        {
            tContent = new T[nSize, nSize];
        }
        public void CreateContentWithInitValue(UInt32 nSize, T TValue)//初始化tContent
        {
            int i,j;
            tContent = new T[nSize, nSize];
            for (i = 0; i < nSize; i++)
                for (j = 0; j < nSize; j++)
                    tContent[i, j] = TValue;
            
        } 

        public int nDataLength
        {
            get
            {
                int nTemp = Convert.ToInt32(System.Runtime.InteropServices.Marshal.SizeOf(typeof(T)) * tContent.GetLength(0) * tContent.GetLength(1));
                return nTemp;
            }
        }
        public int nTotalLength
        {
            get
            {
                if (bExistTag)
                {
                    return nDataLength + 5 + 6 * 2; // 尺寸5字节+边长6字节*2
                }
                else
                {
                    return nDataLength + 6 * 2;//边长10字节
                }
            }
        }

        public TerrainDataError Load(BinaryReader brData, byte[] anBlockTag, bool ExistSizeTag) // 用文件流初始化空类，参数是数组和起始读取的下标,最后一个SizeTag表示是否有数据块长度的信息
        {
            nBlockTag = anBlockTag;
            int nTagLength = anBlockTag.Length;
            byte[] TagRead = brData.ReadBytes(nTagLength);
            if (TagRead[0] != 0x0)
            {
                nSize = (brData.ReadUInt32() - 1) / 2;
                brData.ReadBytes(nXLengthTag.Length);
                UInt32 nXLength = brData.ReadUInt32();
                brData.ReadBytes(nYLengthTag.Length);
                UInt32 nYLength = brData.ReadUInt32();
                UInt32 nDataSize = 0;
                bExistTag = ExistSizeTag;
                if (ExistSizeTag)
                {
                    brData.ReadBytes(nSizeTag.Length);
                    nDataSize = brData.ReadUInt32();
                }
                tContent = new T[nXLength, nYLength];
                switch (typeof(T).ToString())
                {
                    case "System.UInt32":
                        for (int i = 0; i < nXLength; i++)
                        {
                            for (int j = 0; j < nYLength; j++)
                            {
                                tContent[i, j] = (T)(object)brData.ReadUInt32();
                            }
                        }
                        break;
                    case "System.Single":
                        for (int i = 0; i < nXLength; i++)
                        {
                            for (int j = 0; j < nYLength; j++)
                            {
                                tContent[i, j] = (T)(object)brData.ReadSingle();
                            }
                        }
                        break;
                    case "System.Byte":
                        for (int i = 0; i < nXLength; i++)
                        {
                            for (int j = 0; j < nYLength; j++)
                            {
                                tContent[i, j] = (T)(object)brData.ReadByte();
                            }
                        }
                        break;
                    case "System.UInt64":
                        for (int i = 0; i < nXLength; i++)
                        {
                            for (int j = 0; j < nYLength; j++)
                            {
                                tContent[i, j] = (T)(object)brData.ReadUInt64();
                            }
                        }
                        break;
                    default:
                        break;
                }
                return TerrainDataError.None;
            }
            else
            {
                return TerrainDataError.TagZero;
            }
        }
        public string[] ExportText() //假定为正方形，输出文字
        {
            int nLength = tContent.GetLength(0);
            string[] stOutput = new string[nLength];
            for (int i = 0; i < nLength; i++)
            {
                stOutput[i] = "";
                for (int j = 0; j < nLength; j++)
                {
                    /*    string stFormat = "X";
                        switch (typeof(T).ToString())
                        {
                            case "System.Single":
                                stFormat = "F1";
                                break;
                            default:
                                break;
                        }
                        stOutput[i] += tContent[i, j].ToString(stFormat, null) + "  ";
                     * */
                    stOutput[i] += iHexViewer.ToHexString(tContent[nLength-i-1, j]) + "  ";
                }
            }
            return stOutput;
        }
        public Bitmap ExportBMP(int nWidth, int nHeight) // 输出图像
        {
            Bitmap bmpBuffer = new Bitmap(nWidth, nHeight);
            Graphics gBuffer = Graphics.FromImage(bmpBuffer);
            int i, j;
            int nRowsCount = tContent.GetLength(0);
            int nColumnsCount = tContent.GetLength(1);
            float dWidthUnit = ((float)nWidth) / nColumnsCount;
            float dHeightUnit = ((float)nHeight) / nRowsCount;
            for (i = 0; i < nRowsCount; i++)
            {
                for (j = 0; j < nColumnsCount - 1; j++)
                {
                    SolidBrush aBrush = new SolidBrush(iColorize.PointToColor(tContent[nRowsCount-i-1, j]));
                    gBuffer.FillRectangle(aBrush, j * dWidthUnit, i * dHeightUnit, dWidthUnit, dHeightUnit);
                    aBrush.Dispose();
                }
            }
            gBuffer.Dispose();
            return bmpBuffer;
        }

        public void Save(BinaryWriter brData)//向二进制流输出数据
        {
            brData.Write(nBlockTag);
            brData.Write((UInt32)nTotalLength * 2 + 1);
            brData.Write(nXLengthTag);
            brData.Write(tContent.GetLength(0));
            brData.Write(nYLengthTag);
            brData.Write(tContent.GetLength(1));
            if (bExistTag)
            {
                brData.Write(nSizeTag);
                brData.Write((UInt32)nDataLength * 2 + 1);
            }
            switch (typeof(T).ToString())
            {
                case "System.UInt32":
                    for (int i = 0; i < tContent.GetLength(0); i++)
                    {
                        for (int j = 0; j < tContent.GetLength(1); j++)
                        {
                            brData.Write((UInt32)(object)tContent[i, j]);
                        }
                    }
                    break;
                case "System.Single":
                    for (int i = 0; i < tContent.GetLength(0); i++)
                    {
                        for (int j = 0; j < tContent.GetLength(1); j++)
                        {
                            brData.Write((Single)(object)tContent[i, j]);
                        }
                    }
                    break;
                case "System.Byte":
                    for (int i = 0; i < tContent.GetLength(0); i++)
                    {
                        for (int j = 0; j < tContent.GetLength(1); j++)
                        {
                            brData.Write((Byte)(object)tContent[i, j]);
                        }
                    }
                    break;
                case "System.UInt64":
                    for (int i = 0; i < tContent.GetLength(0); i++)
                    {
                        for (int j = 0; j < tContent.GetLength(1); j++)
                        {
                            brData.Write((UInt64)(object)tContent[i, j]);
                        }
                    }
                    break;
                default:
                    break;
            }
        }
    }


}
